package directive

import (
	"errors"
	"strings"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (d Directive) Path() (string, error) {
	var path string

	switch {
	case d.Type() == URL:
		path = d.NamedParameter("Path")

	case d.Type().IsHTTPRequestMethod():
		path = d.NamedParameter("Path")
		if path == "" {
			if d.Parent == nil {
				return "", errors.New(jerr.PathNotFound)
			}
			return d.Parent.Path() // Parent is the URL directive
		}

	default:
		if d.Parent == nil {
			return "", errors.New(jerr.PathNotFound)
		}
		return d.Parent.Path()
	}

	if !strings.HasPrefix(path, "/") {
		return "", errors.New(jerr.IncorrectPath)
	}

	return path, nil
}
