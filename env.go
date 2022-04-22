package main

import (
	"os"
	"strconv"
)

func getBoolEnv(key string) bool {
	val := os.Getenv(key)
	if val == "" {
		return false
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return b
}
