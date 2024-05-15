package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_cors(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		r := httptest.NewRecorder()

		cors(r)

		assert.Equal(t, http.StatusOK, r.Code)
		assert.Len(t, r.Header(), 3)
		assert.Equal(t, "*", r.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE", r.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Browser-UUID", r.Header().Get("Access-Control-Allow-Headers"))
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			cors(nil)
		})
	})
}
