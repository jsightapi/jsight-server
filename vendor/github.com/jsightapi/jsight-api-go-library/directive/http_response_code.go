package directive

import (
	"strconv"
)

func IsHTTPResponseCode(s string) bool {
	if code, err := strconv.Atoi(s); err == nil {
		if s[0] == '0' {
			return false
		}
		if isHTTPResponseCode(code) {
			return true
		}
	}
	return false
}

func isHTTPResponseCode(code int) bool {
	return code >= 100 && code <= 526 // TODO make strict list
}
