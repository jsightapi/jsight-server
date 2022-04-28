package loader

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// Schema compilation works with each node's constraints. This process adjust constraints so that they were in peace
// with each other. For example, if node has "precision" constraint, we have to add to this node "decimal" constraint,
// because only decimal can have precision. Another interesting example: if key-value node doesn't have an "optional"
// constraint, we have to add to parent object "required key" constraint, because this node's key is required in the
// object.

type schemaCompiler struct {
	rootSchema               *schema.Schema
	areKeysOptionalByDefault bool
}

func CompileBasic(rootSchema *schema.Schema, areKeysOptionalByDefault bool) {
	// The root node can be empty, if you specify only ATTRIBUTES without EXAMPLE.
	if rootSchema.RootNode() == nil {
		return
	}

	compile := schemaCompiler{
		rootSchema:               rootSchema,
		areKeysOptionalByDefault: areKeysOptionalByDefault,
	}
	compile.compileNode(rootSchema.RootNode(), 0)
}

func (compile schemaCompiler) compileNode(node schema.Node, indexOfNode int) {
	lex := node.BasisLexEventOfSchemaForNode()
	defer lexeme.CatchLexEventError(lex)

	compile.falseConstraints(node)    // can panic
	compile.orConstraint(node)        // can panic. Must be called before compile.typeConstraint()
	compile.enumConstraint(node)      // can panic. Must be called before compile.typeConstraint()
	compile.precisionConstraint(node) // can panic. Must be called before compile.typeConstraint()
	compile.typeConstraint(node)      // can panic
	if err := compile.allowedConstraintCheck(node); err != nil {
		panic(err)
	}
	compile.anyConstraint(node)                    // can panic
	compile.exclusiveMinimumConstraint(node)       // can panic
	compile.exclusiveMaximumConstraint(node)       // can panic
	compile.optionalConstraints(node, indexOfNode) // can panic
	compile.precisionConstraint(node)              // can panic

	if branchingNode, ok := node.(schema.BranchNode); ok {
		compile.emptyArray(node) // can panic
		for i, child := range branchingNode.Children() {
			compile.compileNode(child, i) // can panic
		}
	}
}

func (schemaCompiler) falseConstraints(node schema.Node) {
	node.ConstraintMap().Filter(func(k constraint.Type, c constraint.Constraint) bool {
		if k == constraint.NullableConstraintType || k == constraint.ConstType {
			if b, ok := c.(constraint.BoolKeeper); ok && !b.Bool() {
				return false
			}
		}
		return true
	})
}

func (schemaCompiler) orConstraint(node schema.Node) {
	if node.Constraint(constraint.OrConstraintType) == nil {
		return
	}

	if node.Constraint(constraint.TypesListConstraintType) == nil {
		panic(errors.ErrLoader) // Links to the schema not found
	}

	// check for a permissible constraints
	n := node.NumberOfConstraints()
	n-- // if node.Constraint(constraint.TypesListConstraintType) != nil - checked above
	if node.Constraint(constraint.OrConstraintType) != nil {
		n--
	}
	if node.Constraint(constraint.OptionalConstraintType) != nil {
		n--
	}
	if node.Constraint(constraint.NullableConstraintType) != nil {
		n--
	}
	if typeConstraint := node.Constraint(constraint.TypeConstraintType); typeConstraint != nil {
		n--
		if t := typeConstraint.(*constraint.TypeConstraint).Bytes().String(); t != `"mixed"` {
			panic(errors.Format(errors.ErrInvalidValueInTheTypeRule, t))
		}
	}
	if n != 0 {
		panic(errors.ErrShouldBeNoOtherRulesInSetWithOr)
	}

	ensureCanUseORConstraint(node)
	node.DeleteConstraint(constraint.OrConstraintType)
}

func ensureCanUseORConstraint(node schema.Node) {
	if branchNode, ok := node.(schema.BranchNode); ok {
		// Since "req.jschema.rules.type.reference 0.2" we didn't allow
		// empty object and arrays as well for the type constraint.
		checkBranchNodeWithOrConstraint(node, branchNode)
	}

	if _, ok := node.(*schema.MixedValueNode); !ok {
		return
	}

	or := node.Constraint(constraint.OrConstraintType).(*constraint.Or) //nolint:errcheck // We are sure about that.
	if or.IsGenerated() {
		return
	}

	ssl := node.Constraint(constraint.TypesListConstraintType).(*constraint.TypesList) //nolint:errcheck // We are sure about that.
	if ssl.HasUserTypes() {
		panic(errors.ErrInvalidChildNodeTogetherWithOrRule)
	}
}

func checkBranchNodeWithOrConstraint(schemaNode schema.Node, jsonNode schema.BranchNode) {
	if jsonNode.Len() != 0 {
		panic(errors.ErrInvalidChildNodeTogetherWithOrRule)
	}

	// Since "req.jschema.rules.or" we didn't allow empty object and arrays for
	// or with at least one user type.
	hasUserTypeInOr := false

	c, ok := schemaNode.Constraint(constraint.TypesListConstraintType).(*constraint.TypesList)
	if !ok {
		return
	}

	for _, n := range c.Names() {
		if n[0] == '@' {
			hasUserTypeInOr = true
			break
		}
	}

	if hasUserTypeInOr {
		panic(errors.ErrInvalidChildNodeTogetherWithOrRule)
	}
}

func (schemaCompiler) enumConstraint(node schema.Node) {
	if node.Constraint(constraint.EnumConstraintType) == nil {
		return
	}

	// check for a permissible constraints
	n := node.NumberOfConstraints()
	n-- // if node.Constraint(constraint.EnumConstraintType) != nil - checked above
	if node.Constraint(constraint.OptionalConstraintType) != nil {
		n--
	}
	if node.Constraint(constraint.ConstType) != nil {
		n--
	}
	if node.Constraint(constraint.NullableConstraintType) != nil {
		n--
	}
	if typeConstraint := node.Constraint(constraint.TypeConstraintType); typeConstraint != nil {
		n--
		if t := typeConstraint.(*constraint.TypeConstraint).Bytes().String(); t != `"enum"` {
			panic(errors.Format(errors.ErrInvalidValueInTheTypeRule, t))
		}
	}
	if n != 0 {
		panic(errors.ErrShouldBeNoOtherRulesInSetWithEnum)
	}
}

func (schemaCompiler) typeConstraint(node schema.Node) { //nolint:gocyclo // todo try to make this more readable
	c := node.Constraint(constraint.TypeConstraintType)
	if c != nil {
		typeConstraint := c.(*constraint.TypeConstraint) //nolint:errcheck // We sure about that.
		val := typeConstraint.Bytes().Unquote()

		if val.IsUserTypeName() {
			n := node.NumberOfConstraints()
			if node.Constraint(constraint.OptionalConstraintType) != nil {
				n--
			}
			if node.Constraint(constraint.NullableConstraintType) != nil {
				n--
			}
			if n != 1 {
				panic(errors.ErrCannotSpecifyOtherRulesWithTypeReference)
			}

			if _, ok := node.(schema.BranchNode); ok {
				// Since "req.jschema.rules.type.reference 0.2" we didn't allow
				// empty object and arrays as well for the type constraint.
				panic(errors.ErrInvalidChildNodeTogetherWithTypeReference)
			}

			if _, ok := node.(*schema.MixedValueNode); ok && !typeConstraint.IsGenerated() {
				// Since "req.jschema.rules.type.reference 0.2" we didn't allow
				// empty object and arrays as well for the type constraint.
				panic(errors.ErrInvalidChildNodeTogetherWithTypeReference)
			}

			c := constraint.NewTypesList(jschema.RuleASTNodeSourceManual)
			c.AddName(val.String(), val.String(), jschema.RuleASTNodeSourceManual)

			node.AddConstraint(c) // can panic: Unable to add constraint
		} else {
			valStr := val.String()

			switch valStr {
			case "mixed":
				typesListConstraint := node.Constraint(constraint.TypesListConstraintType)
				if typesListConstraint == nil {
					panic(errors.ErrNotFoundRuleOr)
				}

				if typesListConstraint.(*constraint.TypesList).Len() < 2 {
					panic(errors.ErrNotFoundRuleOr)
				}
			case "enum":
				if node.Constraint(constraint.EnumConstraintType) == nil {
					panic(errors.ErrNotFoundRuleEnum)
				}
			case "any":
				node.AddConstraint(constraint.NewAny())
			case "decimal":
				if node.Constraint(constraint.PrecisionConstraintType) == nil {
					panic(errors.ErrNotFoundRulePrecision)
				}
			case "email":
				node.AddConstraint(constraint.NewEmail()) // can panic: Unable to add constraint
			case "uri":
				node.AddConstraint(constraint.NewUri())
			case "uuid":
				node.AddConstraint(constraint.NewUuid())
			case "date":
				node.AddConstraint(constraint.NewDate())
			case "datetime":
				node.AddConstraint(constraint.NewDateTime())
			default: // object, array, string, integer, float, boolean or null
				t := json.NewJsonType(val)                         // can panic
				if mixedNode, ok := node.(*schema.MixedNode); ok { // defined json type for mixed node
					mixedNode.SetJsonType(t)
				} else if t != node.Type() { // check json type for non-mixed node
					panic(errors.Format(errors.ErrIncompatibleTypes, t.String()))
				}
			}
			node.SetRealType(valStr)
		}

		node.DeleteConstraint(constraint.TypeConstraintType)
	}
}

func (schemaCompiler) allowedConstraintCheck(node schema.Node) (err error) {
	bannedConstraints := map[constraint.Type][]constraint.Type{
		constraint.EmailConstraintType: {
			constraint.MinLengthConstraintType,
			constraint.MaxLengthConstraintType,
			constraint.RegexConstraintType,
		},

		constraint.UriConstraintType: {
			constraint.MinLengthConstraintType,
			constraint.MaxLengthConstraintType,
			constraint.RegexConstraintType,
		},

		constraint.DateConstraintType: {
			constraint.MinLengthConstraintType,
			constraint.MaxLengthConstraintType,
			constraint.RegexConstraintType,
		},

		constraint.DateTimeConstraintType: {
			constraint.MinLengthConstraintType,
			constraint.MaxLengthConstraintType,
			constraint.RegexConstraintType,
		},

		constraint.UuidConstraintType: {
			constraint.MinLengthConstraintType,
			constraint.MaxLengthConstraintType,
			constraint.RegexConstraintType,
			constraint.RegexConstraintType,
		},

		constraint.AnyConstraintType: {
			constraint.ConstType,
		},
	}

	for t, tt := range bannedConstraints {
		if node.Constraint(t) != nil {
			for _, bt := range tt {
				if node.Constraint(bt) != nil {
					return errors.Format(errors.ErrUnexpectedConstraint, bt.String(), t.String())
				}
			}
		}
	}
	return nil
}

func (schemaCompiler) anyConstraint(node schema.Node) {
	if node.Constraint(constraint.AnyConstraintType) != nil {
		// check for a permissible constraints
		n := node.NumberOfConstraints()
		n-- // AnyConstraintType - checked above
		if node.Constraint(constraint.OptionalConstraintType) != nil {
			n--
		}
		if node.Constraint(constraint.NullableConstraintType) != nil {
			n--
		}
		if node.Constraint(constraint.ConstType) != nil {
			n--
		}
		if n != 0 {
			panic(errors.ErrShouldBeNoOtherRulesInSetWithAny)
		}

		if branchNode, ok := node.(schema.BranchNode); ok {
			if branchNode.Len() != 0 {
				panic(errors.ErrInvalidNestedElementsFoundForTypeAny)
			}
		}
	}
}

func (schemaCompiler) exclusiveMinimumConstraint(node schema.Node) {
	exclusiveMin := node.Constraint(constraint.ExclusiveMinimumConstraintType)
	if exclusiveMin != nil {
		min := node.Constraint(constraint.MinConstraintType)
		if min == nil {
			panic(errors.ErrConstraintMinNotFound)
		}
		if exclusiveMin.(*constraint.ExclusiveMinimum).IsExclusive() {
			min.(*constraint.Min).SetExclusive(true)
		}
		node.DeleteConstraint(constraint.ExclusiveMinimumConstraintType)
	}
}

func (schemaCompiler) exclusiveMaximumConstraint(node schema.Node) {
	exclusiveMax := node.Constraint(constraint.ExclusiveMaximumConstraintType)
	if exclusiveMax != nil {
		max := node.Constraint(constraint.MaxConstraintType)
		if max == nil {
			panic(errors.ErrConstraintMaxNotFound)
		}
		if exclusiveMax.(*constraint.ExclusiveMaximum).IsExclusive() {
			max.(*constraint.Max).SetExclusive(true)
		}
		node.DeleteConstraint(constraint.ExclusiveMaximumConstraintType)
	}
}

func (compile schemaCompiler) optionalConstraints(node schema.Node, indexOfNode int) {
	optional := node.Constraint(constraint.OptionalConstraintType)
	parentNode := node.Parent()
	if optional == nil {
		if objectNode, ok := parentNode.(*schema.ObjectNode); ok {
			if !compile.areKeysOptionalByDefault {
				addRequiredKey(objectNode, objectNode.Key(indexOfNode).Name)
			}
		}
	} else {
		if objectNode, ok := parentNode.(*schema.ObjectNode); ok {
			if !optional.(constraint.BoolKeeper).Bool() {
				addRequiredKey(objectNode, objectNode.Key(indexOfNode).Name)
			}
			node.DeleteConstraint(constraint.OptionalConstraintType)
		} else {
			panic(errors.ErrRuleOptionalAppliesOnlyToObjectProperties)
		}
	}
}

func (schemaCompiler) precisionConstraint(node schema.Node) {
	if node.Constraint(constraint.PrecisionConstraintType) == nil {
		return
	}

	c := node.Constraint(constraint.TypeConstraintType)
	if c == nil {
		return
	}

	t := c.(*constraint.TypeConstraint).Bytes().Unquote().String()
	if t != "decimal" {
		panic(errors.Format(errors.ErrUnexpectedConstraint, "precision", t))
	}
}

func (schemaCompiler) emptyArray(node schema.Node) {
	if arrayNode, ok := node.(*schema.ArrayNode); ok {
		if arrayNode.Len() == 0 {
			zero := json.NewNumberFromUint(0)
			if min := node.Constraint(constraint.MinItemsConstraintType); min != nil {
				if !min.(constraint.ArrayValidator).Value().Equal(zero) {
					panic(errors.ErrIncorrectConstraintValueForEmptyArray)
				}
			}
			if max := node.Constraint(constraint.MaxItemsConstraintType); max != nil {
				if !max.(constraint.ArrayValidator).Value().Equal(zero) {
					panic(errors.ErrIncorrectConstraintValueForEmptyArray)
				}
			}
		}
	}
}
