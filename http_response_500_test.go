package main //nolint:dupl

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_httpResponse500(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		r := httptest.NewRecorder()

		httpResponse500(r, errors.New("fake error"))

		assert.Equal(t, http.StatusInternalServerError, r.Code)
		assert.Len(t, r.Header(), 1)
		assert.Equal(t, "text/plain", r.Header().Get("Content-Type"))
		assert.Equal(t, "fake error", r.Body.String())
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("nil writer", func(t *testing.T) {
			assert.Panics(t, func() {
				httpResponse500(nil, nil)
			})
		})

		t.Run("nil error", func(t *testing.T) {
			assert.Panics(t, func() {
				httpResponse500(httptest.NewRecorder(), nil)
			})
		})
	})
}
