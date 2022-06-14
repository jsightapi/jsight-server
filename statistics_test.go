package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getIP(t *testing.T) {
	const validRemoteAddr = "192.0.2.3:12345"

	createRequest := func(remoteAddr string, headers ...string) *http.Request {
		r, err := http.NewRequest(http.MethodGet, "/", http.NoBody) //nolint:noctx
		require.NoError(t, err)

		r.RemoteAddr = remoteAddr

		if len(headers)%2 != 0 {
			panic("invalid header")
		}

		for i := 0; i < len(headers); i += 2 {
			r.Header.Set(headers[i], headers[i+1])
		}
		return r
	}

	t.Run("positive", func(t *testing.T) {
		cc := map[string]struct {
			given    *http.Request
			expected string
		}{
			"x-real-ip": {
				createRequest(
					validRemoteAddr,
					"X-Real-IP", "192.0.2.1",
					"X-Forwarder-For", "192.0.2.2",
				),
				"192.0.2.1",
			},

			"x-forwarder-for": {
				createRequest(
					validRemoteAddr,
					"X-Forwarder-For", "192.0.2.2",
				),
				"192.0.2.2",
			},

			"remote addr": {
				createRequest(validRemoteAddr),
				"192.0.2.3",
			},

			"cannot split remote addr": {
				createRequest("invalid"),
				"",
			},

			"invalid IP": {
				createRequest("invalid:42"),
				"",
			},
		}

		for n, c := range cc {
			t.Run(n, func(t *testing.T) {
				actual := getIP(c.given)
				assert.Equal(t, c.expected, actual)
			})
		}
	})

	t.Run("negative", func(t *testing.T) {
		assert.Panics(t, func() {
			getIP(nil)
		})
	})
}
