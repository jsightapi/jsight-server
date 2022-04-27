package constraint

import (
	"net/mail"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Email struct{}

var _ Constraint = Email{}

func NewEmail() *Email {
	return &Email{}
}

func (Email) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (Email) Type() Type {
	return EmailConstraintType
}

func (Email) String() string {
	return EmailConstraintType.String()
}

func (Email) Validate(email bytes.Bytes) {
	email = email.Unquote()

	if len(email) == 0 {
		panic(errors.ErrEmptyEmail)
	}

	char := email[0] // first char
	if char == ' ' || char == '<' {
		panic(errors.Format(errors.ErrInvalidEmail, email.String()))
	}

	char = email[len(email)-1] // last char
	if char == ' ' || char == '>' {
		panic(errors.Format(errors.ErrInvalidEmail, email.String()))
	}

	emailStr := email.String()

	_, err := mail.ParseAddress(emailStr)
	if err != nil {
		panic(errors.Format(errors.ErrInvalidEmail, emailStr))
	}
}

func (Email) ASTNode() jschema.RuleASTNode {
	return newEmptyRuleASTNode()
}
