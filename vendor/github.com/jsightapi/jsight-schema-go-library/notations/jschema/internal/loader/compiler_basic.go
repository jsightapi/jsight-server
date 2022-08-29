package loader

import (
	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// schemaCompiler works with each node's constraints. This process adjust constraints
// so that they were in peace with each other. For example, if node has "precision"
// constraint, we have to add to this node "decimal" constraint, because only
// decimal can have precision. Another interesting example: if key-value node doesn't
// have an "optional" constraint, we have to add to parent object "required key"
// constraint, because this node's key is required in the object.
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
	compile.anyConstraint(node) // can panic
	if err := compile.checkPairConstraints(node); err != nil {
		panic(err)
	}
	compile.exclusiveMinimumConstraint(node)       // can panic
	compile.exclusiveMaximumConstraint(node)       // can panic
	compile.optionalConstraints(node, indexOfNode) // can panic

	if branchingNode, ok := node.(schema.BranchNode); ok {
		compile.emptyArray(node) // can panic
		for i, child := range branchingNode.Children() {
			compile.compileNode(child, i) // can panic
		}
	}
}

func (schemaCompiler) falseConstraints(node schema.Node) {
	node.ConstraintMap().Filter(func(k constraint.Type, c constraint.Constraint) bool {
		if k == constraint.NullableConstraintType || k == constraint.ConstConstraintType {
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
	if node.Constraint(constraint.ConstConstraintType) != nil {
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

func (compile schemaCompiler) typeConstraint(node schema.Node) {
	c := node.Constraint(constraint.TypeConstraintType)
	if c == nil {
		return
	}

	typeConstraint := c.(*constraint.TypeConstraint) //nolint:errcheck // We sure about that.
	val := typeConstraint.Bytes().Unquote()

	if val.IsUserTypeName() {
		compile.typeConstraintForUserType(node, typeConstraint, val.String())
	} else {
		compile.typeConstraintForJSONTypes(node, val)
	}

	node.DeleteConstraint(constraint.TypeConstraintType)
}

func (schemaCompiler) typeConstraintForUserType(
	node schema.Node,
	typeConstraint *constraint.TypeConstraint,
	val string,
) {
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
		panic(errors.ErrInvalidChildNodeTogetherWithTypeReference)
	}

	c := constraint.NewTypesList(jschema.RuleASTNodeSourceManual)
	c.AddName(val, val, jschema.RuleASTNodeSourceManual)

	node.AddConstraint(c) // can panic: Unable to add constraint
}

func (schemaCompiler) typeConstraintForJSONTypes(node schema.Node, val bytes.Bytes) {
	valStr := val.String()

	h, ok := jsonTypesHandler[valStr]
	if ok {
		h(node)
	} else {
		t := json.NewJsonType(val)                         // can panic
		if mixedNode, ok := node.(*schema.MixedNode); ok { // defined json type for mixed node
			mixedNode.SetJsonType(t)
		} else if t != node.Type() { // check json type for non-mixed node
			panic(errors.Format(errors.ErrIncompatibleTypes, t.String()))
		}
	}
	if !node.SetRealType(valStr) {
		panic(errors.Format(errors.ErrIncompatibleTypes, valStr))
	}
}

var jsonTypesHandler = map[string]func(node schema.Node){
	"mixed": func(node schema.Node) {
		typesListConstraint := node.Constraint(constraint.TypesListConstraintType)
		if typesListConstraint == nil {
			panic(errors.ErrNotFoundRuleOr)
		}

		if typesListConstraint.(*constraint.TypesList).Len() < 2 {
			panic(errors.ErrNotFoundRuleOr)
		}
	},

	constraint.EnumConstraintType.String(): func(node schema.Node) {
		if node.Constraint(constraint.EnumConstraintType) == nil {
			panic(errors.ErrNotFoundRuleEnum)
		}
	},

	"any": func(node schema.Node) {
		node.AddConstraint(constraint.NewAny())
	},

	"decimal": func(node schema.Node) {
		if node.Constraint(constraint.PrecisionConstraintType) == nil {
			panic(errors.ErrNotFoundRulePrecision)
		}
	},

	"email": func(node schema.Node) {
		node.AddConstraint(constraint.NewEmail()) // can panic: Unable to add constraint
	},

	"uri": func(node schema.Node) {
		node.AddConstraint(constraint.NewUri())
	},

	"uuid": func(node schema.Node) {
		node.AddConstraint(constraint.NewUuid())
	},

	"date": func(node schema.Node) {
		node.AddConstraint(constraint.NewDate())
	},

	"datetime": func(node schema.Node) {
		node.AddConstraint(constraint.NewDateTime())
	},
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
			constraint.ConstConstraintType,
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
	if node.Constraint(constraint.AnyConstraintType) == nil {
		return
	}

	// check for a permissible constraints
	n := node.NumberOfConstraints()
	n-- // AnyConstraintType - checked above
	if node.Constraint(constraint.OptionalConstraintType) != nil {
		n--
	}
	if node.Constraint(constraint.NullableConstraintType) != nil {
		n--
	}
	if node.Constraint(constraint.ConstConstraintType) != nil {
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

// checkPairConstraints checks some constraints which has pairs such as `min` and `max`,
// `minLength` and `maxLength`, etc.
func (compile schemaCompiler) checkPairConstraints(node schema.Node) error {
	checkers := []func(schema.Node) error{
		compile.checkMinAndMax,
		compile.checkMinLengthAndMaxLength,
		compile.checkMinItemsAndMaxItems,
	}

	for _, fn := range checkers {
		if err := fn(node); err != nil {
			return err
		}
	}
	return nil
}

func (schemaCompiler) checkMinAndMax(node schema.Node) error {
	minRaw := node.Constraint(constraint.MinConstraintType)
	maxRaw := node.Constraint(constraint.MaxConstraintType)

	if minRaw == nil || maxRaw == nil {
		return nil
	}

	min := minRaw.(*constraint.Min) //nolint:errcheck // We're sure about this type.
	max := maxRaw.(*constraint.Max) //nolint:errcheck // We're sure about this type.

	if min.Exclusive() || max.Exclusive() {
		if min.Value().GreaterThanOrEqual(max.Value()) {
			return errors.Format(
				errors.ErrValueOfOneConstraintGreaterOrEqualToAnother,
				"min",
				"max",
			)
		}
	} else {
		if min.Value().GreaterThan(max.Value()) {
			return errors.Format(
				errors.ErrValueOfOneConstraintGreaterThanAnother,
				"min",
				"max",
			)
		}
	}
	return nil
}

func (schemaCompiler) checkMinLengthAndMaxLength(node schema.Node) error {
	minLengthRaw := node.Constraint(constraint.MinLengthConstraintType)
	maxLengthRaw := node.Constraint(constraint.MaxLengthConstraintType)

	if minLengthRaw == nil || maxLengthRaw == nil {
		return nil
	}

	minLength := minLengthRaw.(*constraint.MinLength) //nolint:errcheck // We're sure about this type.
	maxLength := maxLengthRaw.(*constraint.MaxLength) //nolint:errcheck // We're sure about this type.

	if minLength.Value() > maxLength.Value() {
		return errors.Format(
			errors.ErrValueOfOneConstraintGreaterThanAnother,
			"minLength",
			"maxLength",
		)
	}
	return nil
}

func (schemaCompiler) checkMinItemsAndMaxItems(node schema.Node) error {
	minItemsRaw := node.Constraint(constraint.MinItemsConstraintType)
	maxItemsRaw := node.Constraint(constraint.MaxItemsConstraintType)

	if minItemsRaw == nil || maxItemsRaw == nil {
		return nil
	}

	minItems := minItemsRaw.(*constraint.MinItems) //nolint:errcheck // We're sure about this type.
	maxItems := maxItemsRaw.(*constraint.MaxItems) //nolint:errcheck // We're sure about this type.

	if minItems.Value() > maxItems.Value() {
		return errors.Format(
			errors.ErrValueOfOneConstraintGreaterThanAnother,
			"minItems",
			"maxItems",
		)
	}
	return nil
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
	objectNode, ok := parentNode.(*schema.ObjectNode)

	if optional == nil {
		if ok && !compile.areKeysOptionalByDefault {
			addRequiredKey(objectNode, objectNode.Key(indexOfNode).Key)
		}
	} else {
		if !ok {
			panic(errors.ErrRuleOptionalAppliesOnlyToObjectProperties)
		}

		if !optional.(constraint.BoolKeeper).Bool() {
			addRequiredKey(objectNode, objectNode.Key(indexOfNode).Key)
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
		panic(errors.Format(errors.ErrUnexpectedConstraint, constraint.PrecisionConstraintType, t))
	}
}

func (schemaCompiler) emptyArray(node schema.Node) {
	arrayNode, ok := node.(*schema.ArrayNode)

	if !ok || arrayNode.Len() != 0 {
		return
	}

	if min := node.Constraint(constraint.MinItemsConstraintType); min != nil {
		if min.(constraint.ArrayValidator).Value() != 0 {
			panic(errors.ErrIncorrectConstraintValueForEmptyArray)
		}
	}
	if max := node.Constraint(constraint.MaxItemsConstraintType); max != nil {
		if max.(constraint.ArrayValidator).Value() != 0 {
			panic(errors.ErrIncorrectConstraintValueForEmptyArray)
		}
	}
}
