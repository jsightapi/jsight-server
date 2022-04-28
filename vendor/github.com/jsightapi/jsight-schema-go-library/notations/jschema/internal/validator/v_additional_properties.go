package validator

import (
	"fmt"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
)

// validator to process additional properties.

type additionalPropertiesValidator struct {
	node_           schema.Node                               // always schema.ObjectNode
	parentValidator validator                                 // always objectValidator
	feedFunc        func(lexeme.LexEvent) ([]validator, bool) // can panic
	schemaType      jschema.SchemaType
	depth           uint
}

// The constructor can return multiple validators because a type can contain an
// "OR" rule.

func newAdditionalPropertiesValidator(
	node schema.Node,
	parentValidator validator,
	c *constraint.AdditionalProperties,
) []validator {
	v := additionalPropertiesValidator{
		node_:           node,            // always schema.ObjectNode
		parentValidator: parentValidator, // always objectValidator
	}

	switch c.Mode() {
	case constraint.AdditionalPropertiesCanBeAny:
		v.feedFunc = v.feedAny

	case constraint.AdditionalPropertiesMustBeSchemaType:
		v.schemaType = c.SchemaType()
		switch v.schemaType {
		case jschema.SchemaTypeObject:
			v.feedFunc = v.feedObject

		case jschema.SchemaTypeArray:
			v.feedFunc = v.feedArray

		default:
			v.feedFunc = v.feedLiteral
		}

	case constraint.AdditionalPropertiesMustBeType:
		schem := parentValidator.(*objectValidator).rootSchema
		return NodeValidatorList(
			schem.Type(c.TypeName().String()).RootNode(), // can panic
			schem,
			parentValidator,
		)

	case constraint.AdditionalPropertiesNotAllowed:
		v.feedFunc = v.feedNotAllowed

	default:
		panic(errors.ErrValidator)
	}

	list := make([]validator, 1)
	list[0] = &v
	return list
}

func (v additionalPropertiesValidator) node() schema.Node {
	return v.node_
}

func (v additionalPropertiesValidator) parent() validator {
	return v.parentValidator
}

func (v *additionalPropertiesValidator) setParent(parent validator) {
	v.parentValidator = parent
}

func (v *additionalPropertiesValidator) feed(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	defer lexeme.CatchLexEventError(jsonLexeme)
	return v.feedFunc(jsonLexeme)
}

func (v *additionalPropertiesValidator) feedAny(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	if jsonLexeme.Type().IsOpening() {
		v.depth++
	} else {
		v.depth--
	}

	if v.depth == 0 {
		return nil, true
	}

	return nil, false
}

func (v *additionalPropertiesValidator) feedObject(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	if jsonLexeme.Type() != lexeme.ObjectBegin {
		panic(errors.ErrUnexpectedLexInObjectValidator)
	}
	v.feedFunc = v.feedAny
	_, b := v.feedFunc(jsonLexeme)
	return nil, b
}

func (v *additionalPropertiesValidator) feedArray(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	if jsonLexeme.Type() != lexeme.ArrayBegin {
		panic(errors.ErrUnexpectedLexInArrayValidator)
	}
	v.feedFunc = v.feedAny
	_, b := v.feedFunc(jsonLexeme)
	return nil, b
}

func (v *additionalPropertiesValidator) feedLiteral(jsonLexeme lexeme.LexEvent) ([]validator, bool) {
	switch jsonLexeme.Type() { //nolint:exhaustive // We will throw a panic in over cases.
	case lexeme.LiteralBegin:
		return nil, false
	case lexeme.LiteralEnd:
		actualType, err := jschema.GuessSchemaType(jsonLexeme.Value())
		if err != nil {
			panic(err)
		}
		if !v.schemaType.IsEqualSoft(actualType) {
			panic(errors.Format(errors.ErrInvalidValueType, actualType, v.schemaType))
		}
		return nil, true
	}
	panic(errors.ErrUnexpectedLexInLiteralValidator)
}

func (*additionalPropertiesValidator) feedNotAllowed(lex lexeme.LexEvent) ([]validator, bool) {
	panic(lexeme.NewLexEventError(
		lex,
		errors.Format(errors.ErrSchemaDoesNotSupportKey, lex.Value().Unquote().String())),
	)
}

func (v additionalPropertiesValidator) log() string {
	return fmt.Sprintf("%s [%p]", v.node_.Type().String(), v.node_)
}
