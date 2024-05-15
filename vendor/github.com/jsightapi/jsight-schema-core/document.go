package schema

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

// Document represents a document.
// It's a concrete data. Data maybe a scalar type or complex type.
//
// Not a thead safe!
//
// Example of the valid documents:
// - "foo"
// - [1, 2, 3]
// - {"foo": "bar"}
type Document interface {
	// NextLexeme returns next lexeme from this document.
	// Might return ParsingError if document isn't valid.
	// Will return io.EOF when no more lexemes are available.
	NextLexeme() (lexeme.LexEvent, error)

	// Len returns length of document in bytes.
	// Might return ParsingError if document isn't valid.
	Len() (uint, error)

	// Check checks that this document is valid.
	// Can return ParsingError if document isn't valid.
	Check() error

	// Content returns the entire contents of the document
	Content() bytes.Bytes
}
