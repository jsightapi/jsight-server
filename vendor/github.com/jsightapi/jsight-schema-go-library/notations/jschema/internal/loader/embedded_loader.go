package loader

import (
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type embeddedLoader interface {
	load(lex lexeme.LexEvent) bool
}
