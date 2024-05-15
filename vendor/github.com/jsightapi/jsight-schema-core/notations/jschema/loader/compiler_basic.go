package loader

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/errs"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

// schemaCompiler works with each node's constraints. This process adjust constraints
// so that they were in peace with each other. For example, if node has "precision"
// constraint, we have to add to this node "decimal" constraint, because only
// decimal can have precision. Another interesting example: if key-value node doesn't
// have an "optional" constraint, we have to add to parent object "required key"
// constraint, because this node's key is required in the object.
type schemaCompiler struct {
	rootSchema               *ischema.ISchema
	areKeysOptionalByDefault bool
}

func CompileBasic(rootSchema *ischema.ISchema, areKeysOptionalByDefault bool) {
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

func (compile schemaCompiler) compileNode(node ischema.Node, indexOfNode int) {
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

	if branchingNode, ok := node.(ischema.BranchNode); ok {
		compile.emptyArray(node) // can panic
		for i, child := range branchingNode.Children() {
			compile.compileNode(child, i) // can panic
		}
	}
}

func (schemaCompiler) falseConstraints(node ischema.Node) {
	node.ConstraintMap().Filter(func(k constraint.Type, c constraint.Constraint) bool {
		if k == constraint.NullableConstraintType || k == constraint.ConstConstraintType {
			if b, ok := c.(constraint.BoolKeeper); ok && !b.Bool() {
				return false
			}
		}
		return true
	})
}

func (schemaCompiler) orConstraint(node ischema.Node) {
	if node.Constraint(constraint.OrConstraintType) == nil {
		return
	}

	if node.Constraint(constraint.TypesListConstraintType) == nil {
		panic(errs.ErrLoader.F()) // Links to the schema not found
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
			panic(errs.ErrInvalidValueInTheTypeRule.F(t))
		}
	}
	if n != 0 {
		panic(errs.ErrShouldBeNoOtherRulesInSetWithOr.F())
	}

	ensureCanUseORConstraint(node)
	node.DeleteConstraint(constraint.OrConstraintType)
}

func ensureCanUseORConstraint(node ischema.Node) {
	if branchNode, ok := node.(ischema.BranchNode); ok {
		// Since "req.jschema.rules.type.reference 0.2" we didn't allow
		// empty object and arrays as well for the type constraint.
		checkBranchNodeWithOrConstraint(node, branchNode)
	}

	if _, ok := node.(*ischema.MixedValueNode); !ok {
		return
	}

	or := node.Constraint(constraint.OrConstraintType).(*constraint.Or)
	if or.IsGenerated() {
		return
	}

	ssl := node.Constraint(constraint.TypesListConstraintType).(*constraint.TypesList)
	if ssl.HasUserTypes() {
		panic(errs.ErrInvalidChildNodeTogetherWithOrRule.F())
	}
}

func checkBranchNodeWithOrConstraint(schemaNode ischema.Node, jsonNode ischema.BranchNode) {
	if jsonNode.Len() != 0 {
		panic(errs.ErrInvalidChildNodeTogetherWithOrRule.F())
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
		panic(errs.ErrInvalidChildNodeTogetherWithOrRule.F())
	}
}

func (schemaCompiler) enumConstraint(node ischema.Node) {
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
			panic(errs.ErrInvalidValueInTheTypeRule.F(t))
		}
	}
	if n != 0 {
		panic(errs.ErrShouldBeNoOtherRulesInSetWithEnum.F())
	}
}

func (compile schemaCompiler) typeConstraint(node ischema.Node) {
	c := node.Constraint(constraint.TypeConstraintType)
	if c == nil {
		return
	}

	typeConstraint := c.(*constraint.TypeConstraint)
	val := typeConstraint.Bytes().Unquote()

	if val.IsUserTypeName() {
		compile.typeConstraintForUserType(node, typeConstraint, val.String())
	} else {
		compile.typeConstraintForJSONTypes(node, val)
	}

	node.DeleteConstraint(constraint.TypeConstraintType)
}

func (schemaCompiler) typeConstraintForUserType(
	node ischema.Node,
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
		panic(errs.ErrCannotSpecifyOtherRulesWithTypeReference.F())
	}

	if _, ok := node.(ischema.BranchNode); ok {
		// Since "req.jschema.rules.type.reference 0.2" we didn't allow
		// empty object and arrays as well for the type constraint.
		panic(errs.ErrInvalidChildNodeTogetherWithTypeReference.F())
	}

	if _, ok := node.(*ischema.MixedValueNode); ok && !typeConstraint.IsGenerated() {
		panic(errs.ErrInvalidChildNodeTogetherWithTypeReference.F())
	}

	c := constraint.NewTypesList(schema.RuleASTNodeSourceManual)
	c.AddName(val, val, schema.RuleASTNodeSourceManual)

	node.AddConstraint(c) // can panic: Unable to add constraint
}

func (schemaCompiler) typeConstraintForJSONTypes(node ischema.Node, val bytes.Bytes) {
	valStr := val.String()

	h, ok := jsonTypesHandler[valStr]
	if ok {
		h(node)
	} else {
		t := json.NewJsonType(val)                          // can panic
		if mixedNode, ok := node.(*ischema.MixedNode); ok { // defined json type for mixed node
			mixedNode.SetJsonType(t)
		} else if t != node.Type() { // check json type for non-mixed node
			panic(errs.ErrIncompatibleTypes.F(t.String()))
		}
	}
	if !node.SetRealType(valStr) {
		panic(errs.ErrIncompatibleTypes.F(valStr))
	}
}

var jsonTypesHandler = map[string]func(node ischema.Node){
	"mixed": func(node ischema.Node) {
		typesListConstraint := node.Constraint(constraint.TypesListConstraintType)
		if typesListConstraint == nil {
			panic(errs.ErrNotFoundRuleOr.F())
		}

		if typesListConstraint.(*constraint.TypesList).Len() < 2 {
			panic(errs.ErrNotFoundRuleOr.F())
		}
	},

	constraint.EnumConstraintType.String(): func(node ischema.Node) {
		if node.Constraint(constraint.EnumConstraintType) == nil {
			panic(errs.ErrNotFoundRuleEnum.F())
		}
	},

	"any": func(node ischema.Node) {
		node.AddConstraint(constraint.NewAny())
	},

	"decimal": func(node ischema.Node) {
		if node.Constraint(constraint.PrecisionConstraintType) == nil {
			panic(errs.ErrNotFoundRulePrecision.F())
		}
	},

	"email": func(node ischema.Node) {
		node.AddConstraint(constraint.NewEmail()) // can panic: Unable to add constraint
	},

	"uri": func(node ischema.Node) {
		node.AddConstraint(constraint.NewUri())
	},

	"uuid": func(node ischema.Node) {
		node.AddConstraint(constraint.NewUuid())
	},

	"date": func(node ischema.Node) {
		node.AddConstraint(constraint.NewDate())
	},

	"datetime": func(node ischema.Node) {
		node.AddConstraint(constraint.NewDateTime())
	},
}

func (schemaCompiler) allowedConstraintCheck(node ischema.Node) (err error) {
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
					return errs.ErrUnexpectedConstraint.F(bt.String(), t.String())
				}
			}
		}
	}
	return nil
}

func (schemaCompiler) anyConstraint(node ischema.Node) {
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
		panic(errs.ErrShouldBeNoOtherRulesInSetWithAny.F())
	}

	if branchNode, ok := node.(ischema.BranchNode); ok {
		if branchNode.Len() != 0 {
			panic(errs.ErrInvalidNestedElementsFoundForTypeAny.F())
		}
	}
}

// checkPairConstraints checks some constraints which has pairs such as `min` and `max`,
// `minLength` and `maxLength`, etc.
func (compile schemaCompiler) checkPairConstraints(node ischema.Node) error {
	checkers := []func(ischema.Node) error{
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

func (schemaCompiler) checkMinAndMax(node ischema.Node) error {
	minRaw := node.Constraint(constraint.MinConstraintType)
	maxRaw := node.Constraint(constraint.MaxConstraintType)

	if minRaw == nil || maxRaw == nil {
		return nil
	}

	min := minRaw.(*constraint.Min)
	max := maxRaw.(*constraint.Max)

	if min.Exclusive() || max.Exclusive() {
		if min.Value().GreaterThanOrEqual(max.Value()) {
			return errs.ErrValueOfOneConstraintGreaterOrEqualToAnother.F("min", "max")
		}
	} else {
		if min.Value().GreaterThan(max.Value()) {
			return errs.ErrValueOfOneConstraintGreaterThanAnother.F("min", "max")
		}
	}
	return nil
}

func (schemaCompiler) checkMinLengthAndMaxLength(node ischema.Node) error {
	minLengthRaw := node.Constraint(constraint.MinLengthConstraintType)
	maxLengthRaw := node.Constraint(constraint.MaxLengthConstraintType)

	if minLengthRaw == nil || maxLengthRaw == nil {
		return nil
	}

	minLength := minLengthRaw.(*constraint.MinLength)
	maxLength := maxLengthRaw.(*constraint.MaxLength)

	if minLength.Value() > maxLength.Value() {
		return errs.ErrValueOfOneConstraintGreaterThanAnother.F(
			"minLength",
			"maxLength",
		)
	}
	return nil
}

func (schemaCompiler) checkMinItemsAndMaxItems(node ischema.Node) error {
	minItemsRaw := node.Constraint(constraint.MinItemsConstraintType)
	maxItemsRaw := node.Constraint(constraint.MaxItemsConstraintType)

	if minItemsRaw == nil || maxItemsRaw == nil {
		return nil
	}

	minItems := minItemsRaw.(*constraint.MinItems)
	maxItems := maxItemsRaw.(*constraint.MaxItems)

	if minItems.Value() > maxItems.Value() {
		return errs.ErrValueOfOneConstraintGreaterThanAnother.F(
			"minItems",
			"maxItems",
		)
	}
	return nil
}

func (schemaCompiler) exclusiveMinimumConstraint(node ischema.Node) {
	exclusiveMin := node.Constraint(constraint.ExclusiveMinimumConstraintType)
	if exclusiveMin != nil {
		min := node.Constraint(constraint.MinConstraintType)
		if min == nil {
			panic(errs.ErrConstraintMinNotFound.F())
		}
		if exclusiveMin.(*constraint.ExclusiveMinimum).IsExclusive() {
			min.(*constraint.Min).SetExclusive(true)
		}
		node.DeleteConstraint(constraint.ExclusiveMinimumConstraintType)
	}
}

func (schemaCompiler) exclusiveMaximumConstraint(node ischema.Node) {
	exclusiveMax := node.Constraint(constraint.ExclusiveMaximumConstraintType)
	if exclusiveMax != nil {
		max := node.Constraint(constraint.MaxConstraintType)
		if max == nil {
			panic(errs.ErrConstraintMaxNotFound.F())
		}
		if exclusiveMax.(*constraint.ExclusiveMaximum).IsExclusive() {
			max.(*constraint.Max).SetExclusive(true)
		}
		node.DeleteConstraint(constraint.ExclusiveMaximumConstraintType)
	}
}

func (compile schemaCompiler) optionalConstraints(node ischema.Node, indexOfNode int) {
	optional := node.Constraint(constraint.OptionalConstraintType)
	parentNode := node.Parent()
	objectNode, ok := parentNode.(*ischema.ObjectNode)

	if optional == nil {
		if ok && !compile.areKeysOptionalByDefault {
			addRequiredKey(objectNode, objectNode.Key(indexOfNode).Key)
		}
	} else {
		if !ok {
			panic(errs.ErrRuleOptionalAppliesOnlyToObjectProperties.F())
		}

		if !optional.(constraint.BoolKeeper).Bool() {
			addRequiredKey(objectNode, objectNode.Key(indexOfNode).Key)
		}
	}
}

func (schemaCompiler) precisionConstraint(node ischema.Node) {
	if node.Constraint(constraint.PrecisionConstraintType) == nil {
		return
	}

	c := node.Constraint(constraint.TypeConstraintType)
	if c == nil {
		return
	}

	t := c.(*constraint.TypeConstraint).Bytes().Unquote().String()
	if t != "decimal" {
		panic(errs.ErrUnexpectedConstraint.F(constraint.PrecisionConstraintType, t))
	}
}

func (schemaCompiler) emptyArray(node ischema.Node) {
	arrayNode, ok := node.(*ischema.ArrayNode)

	if !ok || arrayNode.Len() != 0 {
		return
	}

	if min := node.Constraint(constraint.MinItemsConstraintType); min != nil {
		if min.(constraint.ArrayValidator).Value() != 0 {
			panic(errs.ErrIncorrectConstraintValueForEmptyArray.F())
		}
	}
	if max := node.Constraint(constraint.MaxItemsConstraintType); max != nil {
		if max.(constraint.ArrayValidator).Value() != 0 {
			panic(errs.ErrIncorrectConstraintValueForEmptyArray.F())
		}
	}
}
