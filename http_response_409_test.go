package main //nolint:dupl

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_httpResponse409(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		r := httptest.NewRecorder()

		httpResponse409(r, errors.New("fake error"))

		assert.Equal(t, http.StatusConflict, r.Code)
		assert.Len(t, r.Header(), 1)
		assert.Equal(t, "application/json", r.Header().Get("Content-Type"))
		assert.Equal(t, `{"Status":"Error","Message":"fake error","Line":0,"Index":0}`, r.Body.String())
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("nil writer", func(t *testing.T) {
			assert.Panics(t, func() {
				httpResponse409(nil, nil)
			})
		})

		t.Run("nil error", func(t *testing.T) {
			assert.Panics(t, func() {
				httpResponse409(httptest.NewRecorder(), nil)
			})
		})
	})
}
