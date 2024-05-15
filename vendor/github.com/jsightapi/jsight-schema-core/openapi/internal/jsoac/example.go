package jsoac

import (
	"encoding/json"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

type Example struct {
	value    string
	isString bool
}

var _ json.Marshaler = Example{}
var _ json.Marshaler = &Example{}

// newExample creates an example value for primitive types
func newExample(ex string, isString bool) *Example {
	return &Example{value: ex, isString: isString}
}

func (ex Example) jsonValue() []byte {
	if ex.isString {
		return internal.ToJSONString(ex.value)
	}
	return []byte(ex.value)
}

func (ex Example) MarshalJSON() (b []byte, err error) {
	return ex.jsonValue(), nil
}
