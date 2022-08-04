package checker

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type objectChecker struct{}

func newObjectChecker() objectChecker {
	return objectChecker{}
}

func (objectChecker) Check(nodeLex lexeme.LexEvent) errors.Error {
	if nodeLex.Type() != lexeme.ObjectEnd {
		return lexeme.NewLexEventError(nodeLex, errors.ErrChecker)
	}

	return nil
}
