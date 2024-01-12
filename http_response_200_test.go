package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jsightapi/jsight-api-core/catalog"
	"github.com/stretchr/testify/assert"
)

func Test_httpResponseJDoc200(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("with content", func(t *testing.T) {
			const content = "foobar"
			r := httptest.NewRecorder()

			httpResponseJDoc200(r, []byte(content))

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 2)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, catalog.JDocExchangeVersion, r.Header().Get("X-Jdoc-Exchange-Version"))
			assert.Equal(t, content, r.Body.String())
		})

		t.Run("nil content", func(t *testing.T) {
			r := httptest.NewRecorder()

			httpResponseJDoc200(r, nil)

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 2)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, catalog.JDocExchangeVersion, r.Header().Get("X-Jdoc-Exchange-Version"))
			assert.Equal(t, "", r.Body.String())
		})
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			httpResponseJDoc200(nil, nil)
		})
	})
}

func Test_httpResponseYAML200(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("with content", func(t *testing.T) {
			const content = "foobar"
			r := httptest.NewRecorder()

			httpResponseYAML200(r, []byte(content))

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/yaml; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, content, r.Body.String())
		})

		t.Run("nil content", func(t *testing.T) {
			r := httptest.NewRecorder()

			httpResponseYAML200(r, nil)

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/yaml; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, "", r.Body.String())
		})
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			httpResponseYAML200(nil, nil)
		})
	})
}
