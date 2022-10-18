package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jsightapi/jsight-api-go-library/catalog"
	"github.com/stretchr/testify/assert"
)

func Test_httpResponse200(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("with content", func(t *testing.T) {
			const content = "foobar"
			r := httptest.NewRecorder()

			httpResponse200(r, []byte(content))

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 2)
			assert.Equal(t, "application/json", r.Header().Get("Content-Type"))
			assert.Equal(t, catalog.JDocExchangeVersion, r.Header().Get("X-Jdoc-Exchange-Version"))
			assert.Equal(t, content, r.Body.String())
		})

		t.Run("nil content", func(t *testing.T) {
			r := httptest.NewRecorder()

			httpResponse200(r, nil)

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 2)
			assert.Equal(t, "application/json", r.Header().Get("Content-Type"))
			assert.Equal(t, catalog.JDocExchangeVersion, r.Header().Get("X-Jdoc-Exchange-Version"))
			assert.Equal(t, "", r.Body.String())
		})
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			httpResponse200(nil, nil)
		})
	})
}
