package core

import (
	"fmt"
	"strings"

	"github.com/jsightapi/jsight-api-go-library/catalog"
)

func (core *JApiCore) checkSimilarPaths(pp []PathParameter) error {
	for i := 0; i < len(pp); i++ {
		path := removeLastSegment(pp[i].path)

		if v, ok := core.similarPaths[path]; ok {
			if v != pp[i].parameter {
				return fmt.Errorf("disallow the use of \"similar\" paths: \"/%s/{%s}\", \"/%s\"", path, v, pp[i].path)
			}
		}

		core.similarPaths[path] = pp[i].parameter
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
