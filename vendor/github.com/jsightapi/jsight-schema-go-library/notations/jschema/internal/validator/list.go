package validator

import (
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// Constructor for validator list for a single node. A single node can have multiple validators if, for example, the schema had an "OR" rule.

type validatorListConstructor struct {
	// rootSchema a scheme from which it is possible to receive type by their name.
	rootSchema schema.Schema

	// The parent validator for the newly created validators.
	parent validator

	// addedTypeNames used for excluding recursive addition of type to the list.
	addedTypeNames map[string]struct{}

	// The list of validators for the node.
	list []validator
}

func NodeValidatorList(node schema.Node, rootSchema schema.Schema, parent validator) []validator {
	c := validatorListConstructor{
		rootSchema:     rootSchema,
		parent:         parent,
		addedTypeNames: nil, // optimizing memory allocation
		list:           nil, // optimizing memory allocation
	}
	c.buildList(node)
	return c.list
}

func (c *validatorListConstructor) buildList(node schema.Node) {
	if constr := node.Constraint(constraint.TypesListConstraintType); constr != nil {
		names := constr.(*constraint.TypesList).Names()
		c.appendTypeValidators(names)

		if constr := node.Constraint(constraint.NullableConstraintType); constr != nil {
			c.list = append(c.list, newLiteralValidator(node, c.parent))
		}
	} else {
		c.appendNodeValidators(node)
	}
}

func (c *validatorListConstructor) appendTypeValidators(names []string) {
	if c.list == nil {
		c.addedTypeNames = make(map[string]struct{}, len(names)) // optimizing memory allocation
		c.list = make([]validator, 0, len(names))                // optimizing memory allocation
	}
	for _, name := range names {
		if _, ok := c.addedTypeNames[name]; !ok {
			c.addedTypeNames[name] = struct{}{}
			c.buildList(c.rootSchema.MustType(name).RootNode()) // can panic
		}
	}
}

func (c *validatorListConstructor) appendNodeValidators(node schema.Node) {
	if c.list == nil {
		c.list = make([]validator, 0, 1) // optimizing memory allocation
	}

	var v validator

	if node.Constraint(constraint.AnyConstraintType) != nil && node.Constraint(constraint.ConstConstraintType) == nil {
		v = newAnyNestedStructureValidator(node, c.parent)
	} else {
		switch node.Type() {
		case json.TypeArray:
			v = newArrayValidator(node, c.parent, c.rootSchema)
		case json.TypeObject:
			v = newObjectValidator(node, c.parent, c.rootSchema)
		default:
			v = newLiteralValidator(node, c.parent)
		}
	}

	c.list = append(c.list, v)
}
