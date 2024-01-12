package ischema

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

func newASTNode() schema.ASTNode {
	return schema.ASTNode{
		Rules: &schema.RuleASTNodes{},
	}
}

func astNodeFromNode(n Node) schema.ASTNode {
	an := newASTNode()

	an.TokenType = n.Type().ToTokenType()
	an.SchemaType = getASTNodeSchemaType(n)
	an.Rules = collectASTRules(n.ConstraintMap())
	an.Comment = n.Comment()
	an.InheritedFrom = n.InheritedFrom()

	return an
}

func getASTNodeSchemaType(n Node) string {
	if n.Constraint(constraint.EnumConstraintType) != nil {
		return "enum"
	}

	if n.Constraint(constraint.OrConstraintType) != nil {
		return string(schema.SchemaTypeMixed)
	}

	if c := n.Constraint(constraint.TypeConstraintType); c != nil {
		if tc, ok := c.(*constraint.TypeConstraint); ok {
			return tc.Bytes().Unquote().String()
		}
	}

	if n.Constraint(constraint.PrecisionConstraintType) != nil {
		return string(schema.SchemaTypeDecimal)
	}

	return n.Type().String()
}

func collectASTRules(cc *Constraints) *schema.RuleASTNodes {
	nn := &schema.RuleASTNodes{}

	err := cc.Each(func(k constraint.Type, v constraint.Constraint) error {
		switch k {
		// The `Or` constraint doesn't contain all required values, but they are placed
		// in the `type` constraint.
		case constraint.OrConstraintType:
			types, ok := cc.Get(constraint.TypesListConstraintType)
			if !ok {
				//goland:noinspection GoErrorStringFormat
				return errs.ErrCantCollectRulesTypes.F()
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
