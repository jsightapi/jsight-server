package checker

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

type nodeCheckerListConstructor struct {
	// types the map of used types.
	types map[string]ischema.Type

	// addedTypeNames a set of already added types. Exists for excluding
	// recursive addition of type to the list.
	addedTypeNames map[string]struct{}

	// A rootSchema from which it is possible to receive type by their name.
	rootSchema *ischema.ISchema

	// A list of checkers for the node.
	list []nodeChecker
}

func (l *nodeCheckerListConstructor) buildList(node ischema.Node) {
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

func (l *nodeCheckerListConstructor) appendNodeValidators(node ischema.Node) {
	if l.list == nil {
		l.list = make([]nodeChecker, 0, 1) // optimizing memory allocation
	}

	c, err := newNodeChecker(node)
	if err != nil {
		panic(err)
	}

	l.list = append(l.list, c)
}
