package openapi

import (
	"strings"
)

func concatenateDescription(l, r string) string {
	var sb strings.Builder
	if l != "" {
		sb.WriteString(l)
	}
	if r != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(r)
	}
	return sb.String()
}

//nolint:unused
func escapeTabs(s string) string {
	return strings.ReplaceAll(s, "\t", "\\t")
}

//nolint:unused
func escapeNewLines(s string) string {
	return strings.ReplaceAll(s, "\n", "\\n")
}
