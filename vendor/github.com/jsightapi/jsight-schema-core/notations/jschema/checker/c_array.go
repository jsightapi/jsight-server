package checker

import (
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/kit"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type arrayChecker struct{}

var _ nodeChecker = arrayChecker{}

func newArrayChecker() arrayChecker {
	return arrayChecker{}
}

func (arrayChecker) Check(nodeLex lexeme.LexEvent) kit.Error {
	if nodeLex.Type() != lexeme.ArrayEnd {
		return lexeme.NewError(nodeLex, errs.ErrChecker.F())
	}

	return nil
}
