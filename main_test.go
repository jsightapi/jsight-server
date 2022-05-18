package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_jdocExchangeFile(t *testing.T) { //nolint:funlen
	type testCase struct {
		request  func(*testing.T) *http.Request
		asserter func(*testing.T, *httptest.ResponseRecorder)
	}

	cc := map[string]testCase{
		http.MethodOptions: {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodOptions, "/", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
			},
		},

		"POST, empty request": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
			},
		},

		"POST, with valid schema": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader(`
JSIGHT 0.3

TYPE @cat
	{
		"id": 1
	}
`))
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Equal(t, `{"userTypes":{"@cat":{"schema":{"content":{"properties":{"id":{"jsonType":"number","type":"integer","scalarValue":"1","optional":false}},"jsonType":"object","type":"object","optional":false},"notation":"jsight"}}},"resourceMethods":{},"tags":{},"jdocExchangeFileSchemaVersion":"1.0.0","jsight":"0.3"}`, r.Body.String()) //nolint:lll
			},
		},

		"POST, with invalid schema": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/", strings.NewReader("invalid"))
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusConflict, r.Code)
				assert.Equal(t, `{"Status":"Error","Message":"invalid character 'i' at directive beginning","Line":1,"Index":0}`, r.Body.String()) //nolint:lll
			},
		},
	}

	unhandledMethod := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodTrace,
	}
	for _, m := range unhandledMethod {
		cc[m] = testCase{
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(m, "/", http.NoBody) //nolint:noctx
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusConflict, r.Code)
				assert.Equal(t, "application/json", r.Header().Get("Content-Type"))
				assert.Equal(t, `{"Status":"Error","Message":"HTTP POST request required","Line":0,"Index":0}`, r.Body.String())
			},
		}
	}

	for n, c := range cc {
		t.Run(fmt.Sprintf("%s, without CORS", n), func(t *testing.T) {
			r := httptest.NewRecorder()

			jdocExchangeFile(r, c.request(t))

			c.asserter(t, r)
		})

		t.Run(fmt.Sprintf("%s, with CORS", n), func(t *testing.T) {
			require.NoError(t, os.Setenv("JSIGHT_SERVER_CORS", "true"))
			defer func() {
				require.NoError(t, os.Unsetenv("JSIGHT_SERVER_CORS"))
			}()

			r := httptest.NewRecorder()

			jdocExchangeFile(r, c.request(t))

			c.asserter(t, r)
			assert.Equal(t, "*", r.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE", r.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Browser-UUID", r.Header().Get("Access-Control-Allow-Headers")) //nolint:lll
		})
	}
}
