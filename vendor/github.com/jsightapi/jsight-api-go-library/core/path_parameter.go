package core

import (
	"fmt"
	"strings"

	"github.com/jsightapi/jsight-api-go-library/catalog"
)

type PathParameter struct {
	path      catalog.Path
	parameter string
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	a := strings.Split(path, "/")
	return removeEmptyStrings(a)
}

func removeEmptyStrings(s []string) []string {
	r := make([]string, 0, 3)
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func pathParameters(path string) []PathParameter {
	pp := make([]PathParameter, 0, 3)
	s := splitPath(path)

	for i, segment := range s {
		if segment[0] == '{' && segment[len(segment)-1] == '}' {
			pp = append(pp, PathParameter{
				path:      catalog.Path(strings.Join(s[0:i+1], "/")),
				parameter: segment[1 : len(segment)-1],
			})
		}
	}

	return pp
}

func PathParameters(path string) ([]PathParameter, error) {
	pp := pathParameters(path)

	if hasEmptyPathParameters(pp) {
		return pp, fmt.Errorf(`incorrect empty PATH parameter in "%s"`, path)
	}

	if s := duplicatedPathParameters(pp); s != "" {
		return pp, fmt.Errorf(`the "%s" parameter is duplicated in the path "%s"`, s, path)
	}

	return pp, nil
}

func hasEmptyPathParameters(p []PathParameter) bool {
	for _, pp := range p {
		if pp.parameter == "" {
			return true
		}
	}
	return false
}

// duplicatedPathParameters return duplicated path parameter or empty string if non found
func duplicatedPathParameters(p []PathParameter) string {
	if len(p) <= 1 {
		return ""
	}
	uniq := map[string]struct{}{}
	for _, pp := range p {
		if _, ok := uniq[pp.parameter]; ok {
			return pp.parameter
		}
		uniq[pp.parameter] = struct{}{}
	}
	return ""
}
