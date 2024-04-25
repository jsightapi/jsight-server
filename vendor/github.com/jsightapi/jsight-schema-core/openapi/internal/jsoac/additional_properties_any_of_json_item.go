package jsoac

import (
	"bytes"

	"github.com/jsightapi/jsight-schema-core/openapi/internal"

	"fmt"

	schema "github.com/jsightapi/jsight-schema-core"

	"github.com/jsightapi/jsight-schema-core/errs"

	"strings"
)

type SimpleByteArray []byte

type AdditionalPropertiesAnyJsonItem struct {
	Type                 *string            `json:"type,omitempty"`
	Items                *SimpleByteArray   `json:"items,omitempty"`
	Format               *string            `json:"format,omitempty"`
	Ref                  *string            `json:"$ref,omitempty"`
	Properties           *SimpleByteArray   `json:"properties,omitempty"`
	AdditionalProperties *bool              `json:"additionalProperties,omitempty"`
	Enum                 *[]SimpleByteArray `json:"enum,omitempty"`
	Example              *SimpleByteArray   `json:"example,omitempty"`
}

func makeAdditionalAnyJSONObjects(r schema.RuleASTNode) AdditionalPropertiesAnyJsonItem {
	var s AdditionalPropertiesAnyJsonItem
	var emptyByteObject SimpleByteArray = []byte("{}")

	switch r.Value {
	case internal.StringString, internal.StringInteger, internal.StringBoolean:
		s = AdditionalPropertiesAnyJsonItem{
			Type: internal.StringRef(r.Value),
		}
	case internal.StringFloat:
		s = AdditionalPropertiesAnyJsonItem{
			Type: internal.StringRef(internal.StringNumber),
		}
	case internal.StringArray:
		s = AdditionalPropertiesAnyJsonItem{
			Type:  internal.StringRef(r.Value),
			Items: &emptyByteObject,
		}
	case internal.StringObject:
		s = AdditionalPropertiesAnyJsonItem{
			Type:                 internal.StringRef(internal.StringObject),
			Properties:           &emptyByteObject,
			AdditionalProperties: internal.BoolRef(false),
		}
	case internal.StringNull:
		var nullBytes SimpleByteArray = []byte(internal.StringNull)
		simpleEnum := []SimpleByteArray{
			nullBytes,
		}
		s = AdditionalPropertiesAnyJsonItem{
			Enum:    &simpleEnum,
			Example: &nullBytes,
		}
	case internal.StringDate:
		s = AdditionalPropertiesAnyJsonItem{
			Type:   internal.StringRef(internal.StringString),
			Format: internal.StringRef(internal.StringDate),
		}
	case internal.StringDatetime:
		s = AdditionalPropertiesAnyJsonItem{
			Type:   internal.StringRef(internal.StringString),
			Format: internal.StringRef("date-time"),
		}
	case internal.StringEmail:
		s = AdditionalPropertiesAnyJsonItem{
			Type:   internal.StringRef(internal.StringString),
			Format: internal.StringRef(internal.StringEmail),
		}
	case internal.StringUri:
		s = AdditionalPropertiesAnyJsonItem{
			Type:   internal.StringRef(internal.StringString),
			Format: internal.StringRef(internal.StringUri),
		}
	case internal.StringUuid:
		s = AdditionalPropertiesAnyJsonItem{
			Type:   internal.StringRef(internal.StringString),
			Format: internal.StringRef(internal.StringUuid),
		}
	default:
		if r.Value[0] == '@' {
			s = AdditionalPropertiesAnyJsonItem{
				Ref: internal.StringRef(fmt.Sprintf(`#/components/schemas/%s`, strings.TrimLeft(r.Value, "@"))),
			}
		} else {
			panic(errs.ErrRuntimeFailure.F()) // FIXME: may be: s = AdditionalPropertiesAnyJsonItem{}
		}
	}
	return s
}

func (s SimpleByteArray) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if len(s) == 0 {
		buf.Write([]byte(internal.StringNull))
	} else {
		buf.Write(s)
	}
	return buf.Bytes(), nil
}
