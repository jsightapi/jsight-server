package directive

import (
	"errors"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

func (d Directive) JsonRpcMethodName() (string, error) {
	if d.Type() == Method {
		return d.Parameter("MethodName"), nil
	} else if d.Parent != nil {
		return d.Parent.JsonRpcMethodName()
	}
	return "", errors.New(jerr.JsonRpcMethodNotFound)
}
