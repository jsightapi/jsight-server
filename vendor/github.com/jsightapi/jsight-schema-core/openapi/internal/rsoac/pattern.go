package rsoac

import (
	"encoding/json"
	"strings"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"
)

type Pattern struct {
	value string
}

var _ json.Marshaler = Pattern{}
var _ json.Marshaler = &Pattern{}

// newPattern creates an pattern value for notation regex
func newPattern(ex string) *Pattern {
	return &Pattern{value: ex}
}

func (ex Pattern) jsonValue() []byte {
	value := ex.value
	value = strings.TrimSuffix(value, "/")
	value = strings.TrimPrefix(value, "/")

	return internal.ToJSONString(value)
}

func (ex Pattern) MarshalJSON() (b []byte, err error) {
	return ex.jsonValue(), nil
}
