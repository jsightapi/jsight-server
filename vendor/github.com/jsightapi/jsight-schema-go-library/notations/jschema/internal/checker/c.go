package checker

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type nodeChecker interface {
	check(lexeme.LexEvent) errors.Error
	indentedString(int) string
}
