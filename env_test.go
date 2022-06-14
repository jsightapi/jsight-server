package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_getBoolEnv(t *testing.T) {
	const env = "JSIGHT_GET_BOOL_ENV_TEST"

	cc := map[string]bool{
		"":        false,
		"invalid": false,
		"false":   false,
		"true":    true,
		"True":    true,
		"TRUE":    true,
		"1":       true,
		"t":       true,
	}

	for given, expected := range cc {
		t.Run(given, func(t *testing.T) {
			require.NoError(t, os.Setenv(env, given))

			assert.Equal(t, expected, getBoolEnv(env))
		})
	}

	require.NoError(t, os.Unsetenv(env))
}
