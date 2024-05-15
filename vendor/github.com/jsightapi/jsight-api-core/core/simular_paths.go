package core

import (
	"fmt"
	"strings"

	"github.com/jsightapi/jsight-api-core/jerr"

	"github.com/jsightapi/jsight-api-core/catalog"
)

// checkSimilarPaths Returns an error if similar paths are found. For example: "/cats/{id}" and "/cats/{name}".
func (core *JApiCore) checkSimilarPaths(pp []PathParameter) error {
	for _, p := range pp {
		path := removeLastSegment(p.path)

		if v, ok := core.similarPaths[path]; ok {
			if v != p.parameter {
				var similarPath string
				if path == "" {
					similarPath = fmt.Sprintf("{%s}", v)
				} else {
					similarPath = fmt.Sprintf("%s/{%s}", path, v)
				}
				return fmt.Errorf(jerr.PathsAreSimilar, similarPath, p.path)
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
