package checker

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

type nodeCheckerListConstructor struct {
	// types the map of used types.
	types map[string]schema.Type

	// addedTypeNames a set of already added types. Exists for excluding
	// recursive addition of type to the list.
	addedTypeNames map[string]struct{}

	// A rootSchema from which it is possible to receive type by their name.
	rootSchema *schema.Schema

	// A list of checkers for the node.
	list []nodeChecker
}

func (l *nodeCheckerListConstructor) buildList(node schema.Node) {
	constr := node.Constraint(constraint.TypesListConstraintType)
	if constr != nil {
		names := constr.(*constraint.TypesList).Names()
		l.appendTypeValidators(names)
	} else {
		l.appendNodeValidators(node)
	}
}

func (l *nodeCheckerListConstructor) appendTypeValidators(names []string) {
	if l.list == nil {
		l.addedTypeNames = make(map[string]struct{}, len(names)) // optimizing memory allocation
		l.list = make([]nodeChecker, 0, len(names))              // optimizing memory allocation
	}
	for _, name := range names {
		if _, ok := l.addedTypeNames[name]; !ok {
			l.addedTypeNames[name] = struct{}{}
			l.buildList(getType(name, l.rootSchema, l.types).RootNode()) // can panic
		}
	}
}

func (l *nodeCheckerListConstructor) appendNodeValidators(node schema.Node) {
	if l.list == nil {
		l.list = make([]nodeChecker, 0, 1) // optimizing memory allocation
	}

	var c nodeChecker

	switch node.(type) {
	case *schema.LiteralNode:
		c = newLiteralChecker(node)
	case *schema.ObjectNode:
		c = newObjectChecker(node)
	case *schema.ArrayNode:
		c = newArrayChecker(node)
	case *schema.MixedNode:
		c = newMixedChecker(node)
	default:
		panic(errors.ErrImpossible)
	}

	l.list = append(l.list, c)
}
