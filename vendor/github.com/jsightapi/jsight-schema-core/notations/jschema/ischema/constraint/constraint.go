package constraint

import (
	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type Constraint interface {
	// Type returns the type of constraint.
	Type() Type

	// IsJsonTypeCompatible checks the compatibility of the constraint and json
	// types.
	IsJsonTypeCompatible(json.Type) bool

	// String returns a textual description of the constraint.
	String() string

	// ASTNode returns an AST node for this constraint.
	ASTNode() schema.RuleASTNode
}

type LiteralValidator interface {
	Validate(bytes.Bytes) // Checks the parameter value against the constraint. Panic on an error.
}

type ArrayValidator interface {
	ValidateTheArray(numberOfChildren uint)
	Value() uint
}

type BytesKeeper interface {
	Bytes() bytes.Bytes
}

type BoolKeeper interface {
	Bool() bool
}

// NewConstraintFromRule creates a Constraint from the rule.
// Might return nil.
func NewConstraintFromRule( //nolint:gocyclo // For now it's okay.
	ruleNameLex lexeme.LexEvent,
	ruleValue bytes.Bytes,
	nodeValue bytes.Bytes,
) Constraint {
	str := ruleNameLex.Value().TrimSpaces().Unquote().String()
	switch str {
	case "minLength":
		return NewMinLength(ruleValue)
	case "maxLength":
		return NewMaxLength(ruleValue)
	case "min":
		return NewMin(ruleValue)
	case "max":
		return NewMax(ruleValue)
	case "exclusiveMinimum":
		return NewExclusiveMinimum(ruleValue)
	case "exclusiveMaximum":
		return NewExclusiveMaximum(ruleValue)
	case "type":
		return NewType(ruleValue, schema.RuleASTNodeSourceManual)
	case "precision":
		return NewPrecision(ruleValue)
	case "optional":
		return NewOptional(ruleValue)
	case "minItems":
		return NewMinItems(ruleValue)
	case "maxItems":
		return NewMaxItems(ruleValue)
	case "additionalProperties":
		return NewAdditionalProperties(ruleValue)
	case "nullable":
		return NewNullable(ruleValue)
	case "regex":
		return NewRegex(ruleValue)
	case "const":
		return NewConst(ruleValue, nodeValue)
	}
	panic(lexeme.NewError(ruleNameLex, errs.ErrUnknownRule.F(str)))
}
