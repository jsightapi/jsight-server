package loader

import (
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

func addRequiredKey(node *schema.ObjectNode, key string) {
	requiredKeys := node.Constraint(constraint.RequiredKeysConstraintType)
	if requiredKeys == nil {
		requiredKeys := constraint.NewRequiredKeys()
		requiredKeys.AddKey(key)
		node.AddConstraint(requiredKeys)
	} else {
		requiredKeys.(*constraint.RequiredKeys).AddKey(key)
	}
}
