package loader

import (
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type embeddedLoader interface {
	Load(lex lexeme.LexEvent) bool
}
