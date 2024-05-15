package jsoac

import (
	"encoding/json"
	"regexp"

	schema "github.com/jsightapi/jsight-schema-core"
)

type Description struct {
	value string
}

var _ json.Marshaler = Description{}
var _ json.Marshaler = &Description{}

func newDescription(astNode schema.ASTNode) *Description {
	return newDescriptionFromString(astNode.Comment)
}

func newDescriptionFromString(s string) *Description {
	if len(s) > 0 {
		return &Description{value: Normalize(s)}
	}
	return nil
}

func (ex Description) MarshalJSON() (b []byte, err error) {
	return json.Marshal(ex.value) // JSON quoted string
}

func Normalize(s string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
}
