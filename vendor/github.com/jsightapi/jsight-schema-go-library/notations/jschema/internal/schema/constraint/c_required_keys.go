package constraint

import (
	"strings"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

// RequiredKeys constraint is specific constraint. It cannot be created directly by jSchema language rule.
// It is indirectly influenced by rule "optional" in object's children. All the children keys, that are not marked as
// "optional"=true, are treated as required and go to the RequiredKeys constraint of the parent object.
type RequiredKeys struct {
	keys []string
}

var (
	_ Constraint = RequiredKeys{}
	_ Constraint = (*RequiredKeys)(nil)
)

func NewRequiredKeys() *RequiredKeys {
	return &RequiredKeys{
		keys: make([]string, 0, 10),
	}
}

func (RequiredKeys) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeObject
}

func (RequiredKeys) Type() Type {
	return RequiredKeysConstraintType
}

func (c RequiredKeys) String() string {
	return RequiredKeysConstraintType.String() + ": " + strings.Join(c.keys, ", ")
}

func (c RequiredKeys) Keys() []string {
	return c.keys
}

func (c *RequiredKeys) AddKey(key string) {
	c.keys = append(c.keys, key)
}

func (c RequiredKeys) ASTNode() jschema.RuleASTNode {
	const source = jschema.RuleASTNodeSourceManual

	n := newRuleASTNode(jschema.TokenTypeArray, "", source)
	n.Items = make([]jschema.RuleASTNode, 0, len(c.keys))

	for _, s := range c.keys {
		n.Items = append(n.Items, newRuleASTNode(jschema.TokenTypeString, s, source))
	}

	return n
}
