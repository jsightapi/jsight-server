package directive

import (
	"errors"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (d Directive) HttpMethod() (Enumeration, error) {
	if d.Type().IsHTTPRequestMethod() {
		return d.Type(), nil
	} else if d.Parent != nil {
		return d.Parent.HttpMethod()
	}

	return Get, errors.New(jerr.HttpMethodNotFound)
}
