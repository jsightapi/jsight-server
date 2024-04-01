package jsoac

import (
	"encoding/json"
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
		return toJSONString(ex.value)
	}
	return []byte(ex.value)
}

func (ex Example) MarshalJSON() (b []byte, err error) {
	return ex.jsonValue(), nil
}
