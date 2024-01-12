package checker

import (
	"sort"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

func ValidateLiteralValue(node ischema.Node, jsonValue bytes.Bytes) {
	checkJsonType(node, jsonValue)

	// sorting to make it easier to debug the scheme if there are several errors in it
	m := node.ConstraintMap()
	l := m.Len()
	keys := make([]int, 0, l)

	m.EachSafe(func(k constraint.Type, _ constraint.Constraint) {
		keys = append(keys, int(k))
	})

	sort.Ints(keys)

	if _, ok := m.Get(constraint.NullableConstraintType); ok && jsonValue.String() == "null" {
		return
	}

	for _, k := range keys {
		t := constraint.Type(k)
		c := m.GetValue(t)

		if v, ok := c.(constraint.LiteralValidator); ok {
			v.Validate(jsonValue)
		}
	}
}

var stringBasedTypes = map[constraint.Type]string{
	constraint.EmailConstraintType:    "email",
	constraint.UriConstraintType:      "uri",
	constraint.DateConstraintType:     "date",
	constraint.DateTimeConstraintType: "datetime",
	constraint.UuidConstraintType:     "uuid",
}

func checkJsonType(node ischema.Node, value bytes.Bytes) {
	if node.Constraint(constraint.EnumConstraintType) != nil ||
		node.Constraint(constraint.AnyConstraintType) != nil {
		return
	}

	jsonType := json.Guess(value).LiteralJsonType() // can panic
	schemaType := node.Type()

	if jsonType == json.TypeNull && node.Constraint(constraint.NullableConstraintType) != nil {
		return
	}

	for c, name := range stringBasedTypes {
		if node.Constraint(c) != nil {
			if jsonType == json.TypeString {
				return
			} else {
				panic(errs.ErrInvalidValueType.F(jsonType.String(), name))
			}
		}
	}

	if jsonType == schemaType ||
		(jsonType == json.TypeInteger && schemaType == json.TypeFloat) {
		return
	}

	panic(errs.ErrInvalidValueType.F(jsonType.String(), node.SchemaType()))
}
