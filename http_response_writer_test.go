package main

import (
	"errors"
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
			wr := httpResponseWriter{writer: r}

			wr.jdocJSON([]byte(content))

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 2)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, catalog.JDocExchangeVersion, r.Header().Get("X-Jdoc-Exchange-Version"))
			assert.Equal(t, content, r.Body.String())
		})

		t.Run("nil content", func(t *testing.T) {
			r := httptest.NewRecorder()
			wr := httpResponseWriter{writer: r}

			wr.jdocJSON(nil)

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 2)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, catalog.JDocExchangeVersion, r.Header().Get("X-Jdoc-Exchange-Version"))
			assert.Equal(t, "", r.Body.String())
		})
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			wr := httpResponseWriter{}
			wr.jdocJSON(nil)
		})
	})
}

func Test_httpResponseJSON200(t *testing.T) { //nolint:dupl
	t.Run("positive", func(t *testing.T) {
		t.Run("with content", func(t *testing.T) {
			const content = "foobar"
			r := httptest.NewRecorder()
			wr := httpResponseWriter{writer: r}

			wr.json([]byte(content))

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, content, r.Body.String())
		})

		t.Run("nil content", func(t *testing.T) {
			r := httptest.NewRecorder()
			wr := httpResponseWriter{writer: r}

			wr.json(nil)

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, "", r.Body.String())
		})
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			wr := httpResponseWriter{}
			wr.json(nil)
		})
	})
}

func Test_httpResponseYAML200(t *testing.T) { //nolint:dupl
	t.Run("positive", func(t *testing.T) {
		t.Run("with content", func(t *testing.T) {
			const content = "foobar"
			r := httptest.NewRecorder()
			wr := httpResponseWriter{writer: r}

			wr.yaml([]byte(content))

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/yaml; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, content, r.Body.String())
		})

		t.Run("nil content", func(t *testing.T) {
			r := httptest.NewRecorder()
			wr := httpResponseWriter{writer: r}

			wr.yaml(nil)

			assert.Equal(t, http.StatusOK, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/yaml; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, "", r.Body.String())
		})
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			wr := httpResponseWriter{}
			wr.yaml(nil)
		})
	})
}

func Test_httpResponse409(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("err", func(t *testing.T) {
			r := httptest.NewRecorder()
			wr := httpResponseWriter{writer: r}

			wr.error(errors.New("fake error"))

			assert.Equal(t, http.StatusConflict, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, `{"Status":"Error","Message":"fake error","Line":0,"Index":0}`, r.Body.String())
		})

		t.Run("string", func(t *testing.T) {
			r := httptest.NewRecorder()
			wr := httpResponseWriter{writer: r}

			wr.errorStr("fake error")

			assert.Equal(t, http.StatusConflict, r.Code)
			assert.Len(t, r.Header(), 1)
			assert.Equal(t, "application/json; charset=utf-8", r.Header().Get("Content-Type"))
			assert.Equal(t, `{"Status":"Error","Message":"fake error","Line":0,"Index":0}`, r.Body.String())
		})
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("nil writer", func(t *testing.T) {
			assert.Panics(t, func() {
				wr := httpResponseWriter{}
				wr.error(nil)
			})
		})

		t.Run("nil error", func(t *testing.T) {
			assert.Panics(t, func() {
				r := httptest.NewRecorder()
				wr := httpResponseWriter{writer: r}
				wr.error(nil)
			})
		})
	})
}

func Test_httpResponse500(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		r := httptest.NewRecorder()
		wr := httpResponseWriter{writer: r}

		wr.internalServerError(errors.New("fake error"))

		assert.Equal(t, http.StatusInternalServerError, r.Code)
		assert.Len(t, r.Header(), 1)
		assert.Equal(t, "text/plain", r.Header().Get("Content-Type"))
		assert.Equal(t, "fake error", r.Body.String())
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("nil writer", func(t *testing.T) {
			assert.Panics(t, func() {
				wr := httpResponseWriter{}
				wr.internalServerError(nil)
			})
		})

		t.Run("nil error", func(t *testing.T) {
			assert.Panics(t, func() {
				r := httptest.NewRecorder()
				wr := httpResponseWriter{writer: r}
				wr.internalServerError(nil)
			})
		})
	})
}
