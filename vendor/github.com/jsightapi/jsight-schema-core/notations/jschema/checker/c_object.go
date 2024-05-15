package checker

import (
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type objectChecker struct{}

func newObjectChecker() objectChecker {
	return objectChecker{}
}

func (objectChecker) Check(nodeLex lexeme.LexEvent) kit.Error {
	if nodeLex.Type() != lexeme.ObjectEnd {
		return lexeme.NewError(nodeLex, errs.ErrChecker.F())
	}

	return nil
}
