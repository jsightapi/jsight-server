package enum

import (
	jbytes "github.com/jsightapi/jsight-schema-go-library/bytes"
	jjson "github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type enumItemValue struct {
	value    string
	jsonType jjson.Type
}

func newEnumItem(b jbytes.Bytes) enumItemValue {
	b = b.TrimSpaces()
	t := jjson.Guess(b).JsonType()
	if t == jjson.TypeString {
		b = b.Unquote()
	}
	return enumItemValue{
		value:    b.String(),
		jsonType: t,
	}
}
