package loader

import (
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/lexeme"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema"
	"github.com/jsightapi/jsight-schema-core/notations/jschema/ischema/constraint"
)

type allOfConstraintCompiler struct {
	rootSchema *ischema.ISchema

	// processingTypes a list of schema names that are in the process of compilation
	// (i.e., schemas that contain at least one "allow" rule somewhere inside).
	// Recursive schema processing can occur during compilation.
	processingTypes map[string]struct{}

	// compiledTypes a list of compiled schemas, or schemas that do not need to
	// be compiled (within which the "all Of" rule is not used).
	compiledTypes map[string]struct{}

	foundTypes map[string]ischema.Type
}

// CompileAllOf compile "allOf" rules in root schema, and in all types.
// Adds the necessary properties to objects, removes "allOf" rule.
func CompileAllOf(rootSchema *ischema.ISchema) {
	c := allOfConstraintCompiler{
		rootSchema:      rootSchema,
		processingTypes: make(map[string]struct{}),
		compiledTypes:   make(map[string]struct{}),
		foundTypes:      make(map[string]ischema.Type),
	}

	c.processSchema(rootSchema)

	// In case allow is used only in types (not in the root schema).
	for name := range rootSchema.TypesList() {
		c.processType(name)
	}

	for n, t := range c.foundTypes {
		rootSchema.AddType(n, t)
	}
}

// processSchema searches the schema and processes nodes that contain the "allOf"
// rule.
func (c *allOfConstraintCompiler) processSchema(schem *ischema.ISchema) {
	if node := schem.RootNode(); node != nil {
		c.processNode(node)
	}
}

// processNode recursively searches and processing nodes for the "allOf" rule.
func (c *allOfConstraintCompiler) processNode(node ischema.Node) {
	if allOf := node.Constraint(constraint.AllOfConstraintType); allOf != nil {
		c.extend(node, allOf.(*constraint.AllOf).SchemaNames())
		node.DeleteConstraint(constraint.AllOfConstraintType)
	}

	if branchNode, ok := node.(ischema.BranchNode); ok {
		for _, childNode := range branchNode.Children() {
			c.processNode(childNode)
		}
	}
}

func (c *allOfConstraintCompiler) extend(node ischema.Node, schemaNames []string) {
	defer lexeme.CatchLexEventError(node.BasisLexEventOfSchemaForNode())

	if len(schemaNames) == 0 {
		panic(errs.ErrTypeNameNotFoundInAllOfRule.F())
	}

	for _, name := range schemaNames {
		c.extendWith(node, name)
	}
}

func (c *allOfConstraintCompiler) extendWith(node ischema.Node, name string) {
	lex := node.BasisLexEventOfSchemaForNode()
	defer lexeme.CatchLexEventErrorWithIncorrectUserType(
		lex,
		lex.File().Name(),
	)
	schem := c.processType(name)

	for n, t := range schem.TypesList() {
		c.foundTypes[n] = t
	}

	fromObject, ok := schem.RootNode().(*ischema.ObjectNode)
	if !ok {
		panic(errs.ErrUnacceptableUserTypeInAllOfRule.F(name))
	}

	// It is not obligatory to make a check for casting to type *schema.ObjectNode.
	// The constraint cannot be applied to other types of nodes.
	toObject, ok := node.(*ischema.ObjectNode)
	if !ok {
		panic(errs.ErrUnexpectedConstraint.F(constraint.AllOfConstraintType.String(), node.Type().String()))
	}

	if fromAdditionalProperties := fromObject.Constraint(constraint.AdditionalPropertiesConstraintType); fromAdditionalProperties != nil { //nolint:lll
		fromAdditionalProperties := fromAdditionalProperties.(*constraint.AdditionalProperties)
		if toAdditionalProperties := toObject.Constraint(constraint.AdditionalPropertiesConstraintType); toAdditionalProperties != nil { //nolint:lll
			toAdditionalProperties := toAdditionalProperties.(*constraint.AdditionalProperties)
			if !fromAdditionalProperties.IsEqual(*toAdditionalProperties) {
				panic(errs.ErrConflictAdditionalProperties.F())
			}
		} else {
			toObject.AddConstraint(fromAdditionalProperties)
		}
	}

	for i, childNode := range fromObject.Children() {
		key := fromObject.Key(i)
		cn := childNode.Copy()
		cn.SetInheritedFrom(name)
		toObject.AddChild(key, cn) // can panic ErrDuplicateKeysInSchema
	}

	if requiredKeys := fromObject.Constraint(constraint.RequiredKeysConstraintType); requiredKeys != nil {
		for _, key := range requiredKeys.(*constraint.RequiredKeys).Keys() {
			addRequiredKey(toObject, key)
		}
	}
}

func (c *allOfConstraintCompiler) processType(name string) *ischema.ISchema {
	if _, ok := c.processingTypes[name]; ok {
		panic(errs.ErrUnacceptableRecursionInAllOfRule.F())
	}

	typ := c.rootSchema.MustType(name) // can panic

	if _, ok := c.compiledTypes[name]; ok {
		return typ
	}

	c.processingTypes[name] = struct{}{}
	c.processSchema(typ)
	delete(c.processingTypes, name)

	c.compiledTypes[name] = struct{}{}

	return typ
}
