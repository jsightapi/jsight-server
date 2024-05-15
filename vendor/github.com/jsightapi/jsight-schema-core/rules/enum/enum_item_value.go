package enum

import (
	jbytes "github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/json"
)

type enumItemValue struct {
	value    string
	jsonType json.Type
}

func newEnumItem(b jbytes.Bytes) enumItemValue {
	b = b.TrimSpaces()
	t := json.Guess(b).JsonType()
	if t == json.TypeString {
		b = b.Unquote()
	}
	return enumItemValue{
		value:    b.String(),
		jsonType: t,
	}
}
