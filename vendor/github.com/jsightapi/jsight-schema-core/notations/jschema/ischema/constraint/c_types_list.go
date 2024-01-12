package constraint

import (
	"strings"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/json"
)

type TypesList struct {
	innerTypeNames []string

	// typeNames collection of real type names, used only for building AST nodes.
	typeNames []string

	// elementASTNodes contains an AST node for all items in this constraint
	// This property was added only for build AST node, so it won't affect current
	// logic at all.
	elementASTNodes []schema.RuleASTNode

	source schema.RuleASTNodeSource

	hasUserTypes bool
}

var (
	_ Constraint = TypesList{}
	_ Constraint = (*TypesList)(nil)
)

func NewTypesList(s schema.RuleASTNodeSource) *TypesList {
	return &TypesList{
		innerTypeNames: make([]string, 0, 5),
		source:         s,
	}
}

func (c TypesList) HasUserTypes() bool {
	return c.hasUserTypes
}

func (TypesList) IsJsonTypeCompatible(json.Type) bool {
	return true
}

func (TypesList) Type() Type {
	return TypesListConstraintType
}

func (c TypesList) String() string {
	return TypesListConstraintType.String() + ": " + strings.Join(c.innerTypeNames, ", ")
}

func (c *TypesList) AddName(name, typ string, s schema.RuleASTNodeSource) {
	c.AddNameWithASTNode(name, typ, newRuleASTNode(schema.TokenTypeString, typ, s))
}

func (c *TypesList) AddNameWithASTNode(name, typ string, an schema.RuleASTNode) {
	c.innerTypeNames = append(c.innerTypeNames, name)
	c.typeNames = append(c.typeNames, typ)
	c.elementASTNodes = append(c.elementASTNodes, an)
	c.hasUserTypes = c.hasUserTypes || name[0] == '@'
}

func (c TypesList) Names() []string {
	return c.innerTypeNames
}

func (c TypesList) Len() int {
	return len(c.innerTypeNames)
}

func (c TypesList) ASTNode() schema.RuleASTNode {
	n := newRuleASTNode(schema.TokenTypeArray, "", c.source)
	n.Items = c.elementASTNodes
	return n
}

func (c TypesList) Source() schema.RuleASTNodeSource { return c.source }
