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

type testCase struct {
	request  func(*testing.T) *http.Request
	asserter func(*testing.T, *httptest.ResponseRecorder)
}

func Test_jdocExchangeFile(t *testing.T) {
	cc := map[string]testCase{
		http.MethodOptions: {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodOptions, "/?to=jdoc-2.0", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
			},
		},

		"POST, empty request": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/?to=jdoc-2.0", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
			},
		},

		"POST, with valid schema": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/?to=jdoc-2.0", strings.NewReader(`
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
				assert.Equal(t, `{"tags":{},"userTypes":{"@cat":{"schema":{"content":{"tokenType":"object","type":"object","children":[{"key":"id","tokenType":"number","type":"integer","scalarValue":"1","optional":false}],"optional":false},"example":"{\"id\":1}","notation":"jsight"}}},"interactions":{},"jsight":"0.3","jdocExchangeVersion":"2.0.0"}`, r.Body.String()) //nolint:lll
			},
		},

		"POST, with invalid schema": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/?to=jdoc-2.0", strings.NewReader("invalid"))
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusConflict, r.Code)
				assert.Equal(t, `{"Status":"Error","Message":"invalid character 'i' at directive beginning","Line":1,"Index":0}`, r.Body.String()) //nolint:lll
			},
		},
	}

	appendUnhandledMethod(cc)
	assertAll(t, cc)
}

func Test_openapi(t *testing.T) {
	expectedJSON, err := os.ReadFile("testdata/openapi.json")
	if err != nil {
		require.NoError(t, err)
	}

	expectedYAML, err := os.ReadFile("testdata/openapi.yaml")
	if err != nil {
		require.NoError(t, err)
	}

	cc := map[string]testCase{
		http.MethodOptions: {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodOptions, "/?to=openapi-3.0.3", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
			},
		},

		"POST, default format": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/?to=openapi-3.0.3", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
				assert.Equal(t, expectedJSON, r.Body.Bytes())
			},
		},

		"POST, JSON format": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/?to=openapi-3.0.3&format=json", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
				assert.Equal(t, expectedJSON, r.Body.Bytes())
			},
		},

		"POST, YAML format": {
			func(t *testing.T) *http.Request {
				r, err := http.NewRequest(http.MethodPost, "/?to=openapi-3.0.3&format=yaml", http.NoBody)
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Equal(t, "application/yaml; charset=utf-8", r.Header().Get("Content-Type"))
				assert.Equal(t, expectedYAML, r.Body.Bytes())
			},
		},
	}

	appendUnhandledMethod(cc)
	assertAll(t, cc)
}

func appendUnhandledMethod(cc map[string]testCase) {
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
				r, err := http.NewRequest(m, "/?to=openapi-3.0.3", http.NoBody) //nolint:noctx
				require.NoError(t, err)
				return r
			},
			func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusConflict, r.Code)
				assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
				assert.Equal(t, `{"Status":"Error","Message":"HTTP POST request required","Line":0,"Index":0}`, r.Body.String())
			},
		}
	}
}

func assertAll(t *testing.T, cc map[string]testCase) {
	for n, c := range cc {
		t.Run(fmt.Sprintf("%s, without CORS", n), func(t *testing.T) {
			r := httptest.NewRecorder()

			convertJSight(r, c.request(t))

			c.asserter(t, r)
		})

		t.Run(fmt.Sprintf("%s, with CORS", n), func(t *testing.T) {
			require.NoError(t, os.Setenv("JSIGHT_SERVER_CORS", "true"))
			defer func() {
				require.NoError(t, os.Unsetenv("JSIGHT_SERVER_CORS"))
			}()

			r := httptest.NewRecorder()

			convertJSight(r, c.request(t))

			c.asserter(t, r)
			assert.Equal(t, "*", r.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE", r.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Browser-UUID", r.Header().Get("Access-Control-Allow-Headers")) //nolint:lll
		})
	}
}
