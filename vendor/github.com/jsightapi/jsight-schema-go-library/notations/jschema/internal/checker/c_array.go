package checker

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type arrayChecker struct{}

var _ nodeChecker = arrayChecker{}

func newArrayChecker() arrayChecker {
	return arrayChecker{}
}

func (arrayChecker) Check(nodeLex lexeme.LexEvent) errors.Error {
	if nodeLex.Type() != lexeme.ArrayEnd {
		return lexeme.NewLexEventError(nodeLex, errors.ErrChecker)
	}

	return nil
}
