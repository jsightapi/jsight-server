package schema

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type TokenType = string

const (
	TokenTypeNumber   TokenType = "number"
	TokenTypeString   TokenType = "string"
	TokenTypeBoolean  TokenType = "boolean"
	TokenTypeArray    TokenType = "array"
	TokenTypeObject   TokenType = "object"
	TokenTypeShortcut TokenType = "reference"
	TokenTypeNull     TokenType = "null"
)

type SchemaType string

const (
	SchemaTypeUndefined SchemaType = ""
	SchemaTypeString    SchemaType = "string"
	SchemaTypeInteger   SchemaType = "integer"
	SchemaTypeFloat     SchemaType = "float"
	SchemaTypeDecimal   SchemaType = "decimal"
	SchemaTypeBoolean   SchemaType = "boolean"
	SchemaTypeObject    SchemaType = "object"
	SchemaTypeArray     SchemaType = "array"
	SchemaTypeNull      SchemaType = "null"
	SchemaTypeEmail     SchemaType = "email"
	SchemaTypeURI       SchemaType = "uri"
	SchemaTypeUUID      SchemaType = "uuid"
	SchemaTypeDate      SchemaType = "date"
	SchemaTypeDateTime  SchemaType = "datetime"
	SchemaTypeEnum      SchemaType = "enum"
	SchemaTypeMixed     SchemaType = "mixed"
	SchemaTypeAny       SchemaType = "any"
	SchemaTypeComment   SchemaType = "comment"
)

func IsValidType(s string) bool {
	_, ok := map[string]struct{}{
		string(SchemaTypeString):   {},
		string(SchemaTypeInteger):  {},
		string(SchemaTypeFloat):    {},
		string(SchemaTypeDecimal):  {},
		string(SchemaTypeBoolean):  {},
		string(SchemaTypeObject):   {},
		string(SchemaTypeArray):    {},
		string(SchemaTypeNull):     {},
		string(SchemaTypeEmail):    {},
		string(SchemaTypeURI):      {},
		string(SchemaTypeUUID):     {},
		string(SchemaTypeDate):     {},
		string(SchemaTypeDateTime): {},
		string(SchemaTypeEnum):     {},
		string(SchemaTypeMixed):    {},
		string(SchemaTypeAny):      {},
		string(SchemaTypeComment):  {},
	}[s]
	return ok
}

func (t SchemaType) ToTokenType() string {
	switch t { //nolint:exhaustive // We return an empty string.
	case SchemaTypeObject:
		return "object"
	case SchemaTypeArray:
		return "array"
	case SchemaTypeString:
		return "string"
	case SchemaTypeInteger, SchemaTypeFloat, SchemaTypeDecimal:
		return "number"
	case SchemaTypeBoolean:
		return "boolean"
	case SchemaTypeNull:
		return "null"
	case SchemaTypeMixed:
		return "reference"
	case SchemaTypeComment:
		return "annotation"
	}
	return ""
}

func (t SchemaType) IsScalar() bool {
	return t.IsOneOf(
		SchemaTypeString,
		SchemaTypeInteger,
		SchemaTypeFloat,
		SchemaTypeDecimal,
		SchemaTypeBoolean,
		SchemaTypeNull,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
	)
}

// IsOneOf return true if current schema is one of specified.
func (t SchemaType) IsOneOf(tt ...SchemaType) bool {
	if t == SchemaTypeUndefined {
		return false
	}

	for _, x := range tt {
		if t == x {
			return true
		}
	}
	return false
}

// IsEqualSoft compare two types with next assumptions%
// - Decimal is the same as float;
// - Email, URI, UUID, Date, and DateTime are the same as string;
// - Enum, Mixed and Any are the same as any other type.
func (t SchemaType) IsEqualSoft(x SchemaType) bool {
	// Fast path.
	if t == x {
		return t != SchemaTypeUndefined
	}

	// Slow path.
	for _, r := range schemaTypeComparisonMap[t] {
		if x == r {
			return true
		}
	}
	return false
}

var schemaTypeComparisonMap = map[SchemaType][]SchemaType{
	SchemaTypeUndefined: {},
	SchemaTypeString: {
		SchemaTypeString,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeInteger: {
		SchemaTypeInteger,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeFloat: {
		SchemaTypeFloat,
		SchemaTypeDecimal,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeDecimal: {
		SchemaTypeFloat,
		SchemaTypeDecimal,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeBoolean: {
		SchemaTypeBoolean,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeObject: {
		SchemaTypeObject,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeArray: {
		SchemaTypeArray,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeNull: {
		SchemaTypeNull,
		SchemaTypeArray,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeEmail: {
		SchemaTypeString,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeURI: {
		SchemaTypeString,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeUUID: {
		SchemaTypeString,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeDate: {
		SchemaTypeString,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeDateTime: {
		SchemaTypeString,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeEnum: {
		SchemaTypeString,
		SchemaTypeInteger,
		SchemaTypeFloat,
		SchemaTypeDecimal,
		SchemaTypeBoolean,
		SchemaTypeObject,
		SchemaTypeArray,
		SchemaTypeNull,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeMixed: {
		SchemaTypeString,
		SchemaTypeInteger,
		SchemaTypeFloat,
		SchemaTypeDecimal,
		SchemaTypeBoolean,
		SchemaTypeObject,
		SchemaTypeArray,
		SchemaTypeNull,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
	SchemaTypeAny: {
		SchemaTypeString,
		SchemaTypeInteger,
		SchemaTypeFloat,
		SchemaTypeDecimal,
		SchemaTypeBoolean,
		SchemaTypeObject,
		SchemaTypeArray,
		SchemaTypeNull,
		SchemaTypeEmail,
		SchemaTypeURI,
		SchemaTypeUUID,
		SchemaTypeDate,
		SchemaTypeDateTime,
		SchemaTypeEnum,
		SchemaTypeMixed,
		SchemaTypeAny,
	},
}

func GuessSchemaType(b []byte) (SchemaType, error) {
	return (&typeGuesser{data: b}).Guess()
}

type typeGuesser struct {
	number *json.Number
	data   []byte
}

func (g *typeGuesser) Guess() (SchemaType, error) {
	m := map[SchemaType]func() bool{
		SchemaTypeString:  g.isString,
		SchemaTypeInteger: g.isInteger,
		SchemaTypeFloat:   g.isFloat,
		SchemaTypeBoolean: g.isBoolean,
		SchemaTypeObject:  g.isObject,
		SchemaTypeArray:   g.isArray,
		SchemaTypeNull:    g.isNull,
	}

	for t, fn := range m {
		if fn() {
			return t, nil
		}
	}
	return SchemaTypeUndefined, errs.ErrUnableToDetermineTheTypeOfJsonValue.F()
}

func (g *typeGuesser) isString() bool {
	length := len(g.data)
	return length >= 2 && g.data[0] == '"' && g.data[length-1] == '"'
}

func (g *typeGuesser) isInteger() bool {
	dot := false
	exp := false
	for _, c := range g.data {
		switch c {
		case '.':
			dot = true
		case 'e', 'E':
			exp = true
		}
	}
	if dot && !exp {
		return false
	}

	n, err := g.parseNumber()
	if err != nil {
		return false
	}

	if n.LengthOfFractionalPart() != 0 {
		return false
	}

	return true
}

func (g *typeGuesser) isFloat() bool {
	dot := false
	exp := false
	for _, c := range g.data {
		switch c {
		case '.':
			dot = true
		case 'e', 'E':
			exp = true
		}
	}
	if dot && !exp {
		return true
	}

	n, err := g.parseNumber()
	if err != nil {
		return false
	}

	if n.LengthOfFractionalPart() != 0 {
		return true
	}

	return false
}

func (g *typeGuesser) isBoolean() bool {
	str := string(g.data)
	return str == "true" || str == "false"
}

func (g *typeGuesser) isObject() bool {
	return string(g.data) == "{"
}

func (g *typeGuesser) isArray() bool {
	return string(g.data) == "["
}

func (g *typeGuesser) isNull() bool {
	return string(g.data) == "null"
}

func (g *typeGuesser) parseNumber() (*json.Number, error) {
	if g.number == nil {
		n, err := json.NewNumber(bytes.NewBytes(g.data))
		if err != nil {
			return nil, err
		}
		g.number = n
	}
	return g.number, nil
}
