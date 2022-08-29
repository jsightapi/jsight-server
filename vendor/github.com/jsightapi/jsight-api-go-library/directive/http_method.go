package directive

import (
	"errors"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (d Directive) HTTPMethod() (Enumeration, error) {
	if d.Type().IsHTTPRequestMethod() {
		return d.Type(), nil
	} else if d.Parent != nil {
		return d.Parent.HTTPMethod()
	}

	return Get, errors.New(jerr.HTTPMethodNotFound)
}
