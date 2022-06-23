package schema

import (
	"errors"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

func newASTNode() jschema.ASTNode {
	return jschema.ASTNode{
		Rules: &jschema.RuleASTNodes{},
	}
}

func astNodeFromNode(n Node) jschema.ASTNode {
	an := newASTNode()

	an.JSONType = n.Type().ToTokenType()
	an.SchemaType = getASTNodeSchemaType(n)
	an.Rules = collectASTRules(n.ConstraintMap())
	an.Comment = n.Comment()

	return an
}

func getASTNodeSchemaType(n Node) string {
	if n.Constraint(constraint.EnumConstraintType) != nil {
		return "enum"
	}

	if n.Constraint(constraint.OrConstraintType) != nil {
		return jschema.JSONTypeMixed
	}

	if c := n.Constraint(constraint.TypeConstraintType); c != nil {
		if tc, ok := c.(*constraint.TypeConstraint); ok {
			return tc.Bytes().Unquote().String()
		}
	}

	if n.Constraint(constraint.PrecisionConstraintType) != nil {
		return "decimal"
	}

	return n.Type().String()
}

func collectASTRules(cc *Constraints) *jschema.RuleASTNodes {
	nn := &jschema.RuleASTNodes{}

	err := cc.Each(func(k constraint.Type, v constraint.Constraint) error {
		switch k {
		// The `Or` constraint doesn't contain all required values, but they are placed
		// in the `type` constraint.
		case constraint.OrConstraintType:
			types, ok := cc.Get(constraint.TypesListConstraintType)
			if !ok {
				//goland:noinspection GoErrorStringFormat
				return errors.New(`Can't collect rules: "types" constraint is required with "or"" constraint`) //nolint:stylecheck // It's okay.
			}

			nn.Set(constraint.OrConstraintType.String(), types.ASTNode())

		case constraint.TypesListConstraintType:
			// do nothing

		default:
			nn.Set(k.String(), v.ASTNode())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return nn
}
