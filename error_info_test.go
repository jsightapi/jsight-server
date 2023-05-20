package main

import (
	"errors"
	"testing"

	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-schema-core/fs"
	"github.com/stretchr/testify/assert"
)

func Test_newErrorInfo(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given    error
			expected errorInfo
		}{
			"ordinal error": {
				errors.New("fake error"),
				errorInfo{
					Status:  "Error",
					Message: "fake error",
				},
			},

			"JAPI error": {
				jerr.NewJApiError("fake error", fs.NewFile("foo", []byte("123")), 2),
				errorInfo{
					Status:  "Error",
					Message: "fake error",
					Line:    1,
					Index:   2,
				},
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				actual := newErrorInfo(c.given)
				assert.Equal(t, c.expected, actual)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			newErrorInfo(nil)
		})
	})
}
