package core

import (
	"fmt"
	"strings"

	"github.com/jsightapi/jsight-api-go-library/catalog"
)

func (core *JApiCore) checkSimilarPaths(pp []PathParameter) error {
	for _, p := range pp {
		path := removeLastSegment(p.path)

		if v, ok := core.similarPaths[path]; ok {
			if v != p.parameter {
				return fmt.Errorf("disallow the use of \"similar\" paths: \"/%s/{%s}\", \"/%s\"", path, v, p.path)
			}
		}

		core.similarPaths[path] = p.parameter
	}
	return nil
}

func removeLastSegment(p catalog.Path) string {
	ss := splitPath(string(p))
	if len(ss) != 0 {
		ss = ss[:len(ss)-1]
	}
	return strings.Join(ss, "/")
}
