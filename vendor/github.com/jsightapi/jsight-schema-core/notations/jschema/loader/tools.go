package loader

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

func addRequiredKey(node *ischema.ObjectNode, key string) {
	requiredKeys := node.Constraint(constraint.RequiredKeysConstraintType)
	if requiredKeys == nil {
		requiredKeys := constraint.NewRequiredKeys()
		requiredKeys.AddKey(key)
		node.AddConstraint(requiredKeys)
	} else {
		requiredKeys.(*constraint.RequiredKeys).AddKey(key)
	}
}
