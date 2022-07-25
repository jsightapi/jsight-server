package checker

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// Checks the SAMPLE SCHEMA and all TYPES for compliance with all RULES.

type checkSchema struct {
	rootSchema *schema.Schema

	// foundTypeNames the names of the type encountered during checking. Are used
	// to control recursion.
	foundTypeNames map[string]struct{}

	// allowedJsonTypes the list of available json-types from types.
	allowedJsonTypes map[json.Type]struct{}
}

func CheckRootSchema(rootSchema *schema.Schema) {
	c := checkSchema{
		rootSchema:       rootSchema,
		foundTypeNames:   make(map[string]struct{}, 10),
		allowedJsonTypes: make(map[json.Type]struct{}, 10),
	}

	if rootSchema.RootNode() != nil { // the root schema may contain no nodes
		c.checkNode(rootSchema.RootNode(), rootSchema.TypesList())
	}

	for name, typ := range rootSchema.TypesList() {
		c.checkType(name, typ, rootSchema.TypesList())
	}
}

func (c *checkSchema) checkType(name string, typ schema.Type, ss map[string]schema.Type) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		// Return an error with the full set of bytes of the root schema.
		if documentError, ok := r.(errors.DocumentError); ok {
			documentError.SetFile(typ.RootFile())
			documentError.SetIndex(documentError.Index() + typ.Begin())
			documentError.SetIncorrectUserType(name)
			panic(documentError)
		}

		panic(r)
	}()

	c.checkNode(typ.Schema().RootNode(), ss)
}

func (c *checkSchema) checkerList(node schema.Node, ss map[string]schema.Type) []nodeChecker {
	l := nodeCheckerListConstructor{
		rootSchema: c.rootSchema,
		types:      ss,
	}
	l.buildList(node)
	return l.list
}

func (c checkSchema) checkNode(node schema.Node, ss map[string]schema.Type) {
	defer lexeme.CatchLexEventError(node.BasisLexEventOfSchemaForNode())
	switch node := node.(type) {
	case *schema.LiteralNode:
		// c.log.Default(node.IndentedNodeString(0))
		c.checkCompatibilityOfConstraints(node)
		c.checkLinksOfNode(node, ss) // can panic
		c.checkLiteralNode(node, ss)
	case *schema.ArrayNode:
		c.checkCompatibilityOfConstraints(node)
		c.checkLinksOfNode(node, ss) // can panic
		c.checkArrayItems(node)
	case *schema.ObjectNode:
		c.checkCompatibilityOfConstraints(node)
		c.checkLinksOfNode(node, ss) // can panic
		c.ensureShortcutKeysAreValid(node)
	case *schema.MixedNode:
		c.checkCompatibilityOfConstraints(node)
		c.checkLinksOfNode(node, ss) // can panic
	case *schema.MixedValueNode:
		c.checkCompatibilityOfConstraints(node)
		c.checkLinksOfNode(node, ss) // can panic
	default:
		panic(errors.ErrImpossible)
	}

	if branchingNode, ok := node.(schema.BranchNode); ok {
		for _, child := range branchingNode.Children() {
			c.checkNode(child, ss) // can panic
		}
	}
}

func (c checkSchema) checkLiteralNode(node schema.Node, ss map[string]schema.Type) {
	checkerList := c.checkerList(node, ss)
	errorsCount := 0
	var err errors.Error

	for _, checker := range checkerList {
		err = checker.check(node.BasisLexEventOfSchemaForNode())
		if err != nil {
			errorsCount++
		}
	}

	if errorsCount == len(checkerList) {
		if len(checkerList) == 1 {
			panic(err)
		} else {
			panic(lexeme.NewLexEventError(node.BasisLexEventOfSchemaForNode(), errors.ErrOrRuleSetValidation))
		}
	}
}

// Checks for array elements. Including recursively for types. Or if the array
// type is "any".
func (c checkSchema) checkArrayItems(node schema.Node) {
	arrayNode := node.(*schema.ArrayNode) //nolint:errcheck // We're sure about this type.

	if arrayNode.Len() != 0 {
		return
	}

	if arrayNode.Constraint(constraint.AnyConstraintType) != nil {
		return
	}

	if typesList := arrayNode.Constraint(constraint.TypesListConstraintType); typesList != nil {
		for _, name := range typesList.(*constraint.TypesList).Names() {
			typeRootNode := c.rootSchema.Type(name).RootNode() // can panic

			if arrayNode, ok := typeRootNode.(*schema.ArrayNode); ok {
				c.checkArrayItems(arrayNode)
			}
		}
	}
}

// check all constraints for compatibility with the json-type of the node
func (checkSchema) checkCompatibilityOfConstraints(node schema.Node) {
	_, isMixed := node.(*schema.MixedNode)
	_, isMixedValue := node.(*schema.MixedValueNode)

	err := node.ConstraintMap().Each(func(k constraint.Type, v constraint.Constraint) error {
		if !v.IsJsonTypeCompatible(node.Type()) && !isMixed && !isMixedValue {
			return errors.Format(errors.ErrUnexpectedConstraint, v.Type().String(), node.RealType())
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func (c *checkSchema) checkLinksOfNode(node schema.Node, ss map[string]schema.Type) {
	if node.Constraint(constraint.TypesListConstraintType) == nil {
		return // to optimize memory allocation
	}

	for k := range c.foundTypeNames {
		delete(c.foundTypeNames, k)
	}
	for k := range c.allowedJsonTypes {
		delete(c.allowedJsonTypes, k)
	}

	c.collectAllowedJsonTypes(node, ss)
	if _, ok := c.allowedJsonTypes[node.Type()]; !ok {
		panic(errors.ErrIncorrectUserType)
	}
}

func (c *checkSchema) ensureShortcutKeysAreValid(node *schema.ObjectNode) {
	var lex lexeme.LexEvent

	defer func() {
		r := recover()
		if r == nil {
			return
		}

		err, ok := r.(errors.Errorf)
		if !ok {
			panic(r)
		}

		if err.Code() != errors.ErrTypeNotFound {
			panic(r)
		}

		panic(lexeme.NewLexEventError(lex, err))
	}()

	for _, v := range node.Keys().Data {
		if v.IsShortcut {
			lex = v.Lex
			c.rootSchema.Type(v.Key) // can panic
		}
	}
}

func (c *checkSchema) collectAllowedJsonTypes(node schema.Node, ss map[string]schema.Type) {
	if _, ok := node.(*schema.MixedValueNode); ok {
		// This node can be anything.
		for _, t := range json.AllTypes {
			c.allowedJsonTypes[t] = struct{}{}
		}

		// Check all user types are defined.
		if typesConstraint := node.Constraint(constraint.TypesListConstraintType); typesConstraint != nil {
			for _, typeName := range typesConstraint.(*constraint.TypesList).Names() {
				c.rootSchema.Type(typeName) // can panic
			}
		}
		return
	}

	typesConstraint := node.Constraint(constraint.TypesListConstraintType)

	if typesConstraint == nil {
		c.allowedJsonTypes[node.Type()] = struct{}{}
		return
	}

	for _, typeName := range typesConstraint.(*constraint.TypesList).Names() {
		if _, ok := c.foundTypeNames[typeName]; ok {
			panic(errors.Format(errors.ErrImpossibleToDetermineTheJsonTypeDueToRecursion, typeName))
		}
		c.foundTypeNames[typeName] = struct{}{}
		c.collectAllowedJsonTypes(getType(typeName, c.rootSchema, ss).RootNode(), ss) // can panic
	}
}

func getType(n string, rootSchema *schema.Schema, ss map[string]schema.Type) (ret *schema.Schema) {
	getFromRoot := func() *schema.Schema {
		return rootSchema.Type(n)
	}

	getFromMap := func() *schema.Schema {
		s, ok := ss[n]
		if !ok {
			panic(errors.Format(errors.ErrTypeNotFound, n))
		}
		return s.Schema()
	}

	main := getFromRoot
	alternative := getFromMap
	if len(n) > 0 && n[0] == '#' {
		main = getFromMap
		alternative = getFromRoot
	}

	defer func() {
		if r := recover(); r == nil {
			return
		}

		ret = alternative()
	}()
	return main()
}
