package jschema

import (
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

func UserTypeNamesFromEachTypeConstraint(node ischema.Node) []string {
	res := make([]string, 0, 3)
	res = append(res, UserTypeNamesFromTypeConstraint(node)...)
	res = append(res, UserTypeNamesFromTypesListConstraint(node)...)
	return res
}

func UserTypeNamesFromTypeConstraint(node ischema.Node) []string {
	cnstr := node.Constraint(constraint.TypeConstraintType)
	if cnstr == nil {
		return nil
	}

	typ, ok := cnstr.(*constraint.TypeConstraint)
	if !ok {
		return nil
	}

	name := typ.Bytes().Unquote().String()
	if name[0] == '@' {
		return []string{name}
	}

	return nil
}

func UserTypeNamesFromTypesListConstraint(node ischema.Node) []string {
	cnstr := node.Constraint(constraint.TypesListConstraintType)
	if cnstr == nil {
		return nil
	}

	list, ok := cnstr.(*constraint.TypesList)
	if !ok {
		return nil
	}

	res := make([]string, 0, list.Len())

	for _, name := range list.Names() {
		if name[0] == '@' {
			res = append(res, name)
		}
	}

	return res
}
