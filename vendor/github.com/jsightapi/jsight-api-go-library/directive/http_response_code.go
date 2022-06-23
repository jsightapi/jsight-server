package directive

import (
	"strconv"
)

func IsHTTPResponseCode(s string) bool {
	code, err := strconv.Atoi(s)
	if err != nil {
		return false
	}

	if s[0] == '0' {
		return false
	}

	return isHTTPResponseCode(code)
}

func isHTTPResponseCode(code int) bool {
	return code >= 100 && code <= 526 // TODO make strict list
}
