package catalog

import (
	"regexp"
	"strings"
)

func Annotation(s string) string {
	return annotationReplacer.ReplaceAllString(strings.TrimSpace(s), " ")
}

var annotationReplacer = regexp.MustCompile(`\s+`)
