package scanner

import (
	"fmt"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type LexemeType uint8

const (
	// Keyword represent a name of Directive, i.e. URL, GET, Path, 200, etc.
	Keyword LexemeType = iota

	// Parameter represent a parameter for directive (regexp/jsight for TYPE)
	Parameter

	// Annotation represent user's annotation to directive in a free-text form.
	Annotation

	// Schema represent a jSchema inside directive's body (Body, TYPE, 200, etc).
	Schema

	// Json represents a JSON inside directive's body (CONFIG).
	Json

	// Text represents a text inside directive Description body.
	Text

	// ContextExplicitOpening represent explicitly opened context, so that it can
	// be later explicitly closed.
	ContextExplicitOpening

	// ContextExplicitClosing represents explicitly closed context.
	ContextExplicitClosing

	// Enum represents an enum body.
	Enum
)

func (t LexemeType) String() string {
	if s, ok := lexemeTypeStringMap[t]; ok {
		return s
	}
	return "unknown-lexeme-type"
}

var lexemeTypeStringMap = map[LexemeType]string{
	Keyword:                "keyword",
	Parameter:              "property",
	Annotation:             "annotation",
	Schema:                 "schema",
	Json:                   "json",
	Text:                   "text",
	ContextExplicitOpening: "context-opening",
	ContextExplicitClosing: "context-closing",
}

type Lexeme struct {
	file  *fs.File    // File containing the contents of the json and the file name
	type_ LexemeType  // Type of found lexeme
	begin bytes.Index // bytes.Index of the start character of the found lexeme in the byte slice
	end   bytes.Index // bytes.Index of the end character of the found lexeme in the byte slice
}

func NewLexeme(type_ LexemeType, begin bytes.Index, end bytes.Index, file *fs.File) *Lexeme {
	return &Lexeme{
		file:  file,
		type_: type_,
		begin: begin,
		end:   end,
	}
}

func (lex Lexeme) Value() bytes.Bytes {
	return lex.file.Content().Slice(lex.begin, lex.end)
}

func (lex Lexeme) File() *fs.File {
	return lex.file
}

func (lex Lexeme) Type() LexemeType {
	return lex.type_
}

func (lex Lexeme) Begin() bytes.Index {
	return lex.begin
}

func (lex Lexeme) End() bytes.Index {
	return lex.end
}

func (lex Lexeme) String() string {
	return fmt.Sprintf("%s [%d:%d]", lex.type_.String(), lex.begin, lex.end)
}
