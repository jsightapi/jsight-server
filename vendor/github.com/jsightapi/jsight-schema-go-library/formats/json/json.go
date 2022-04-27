package json

import (
	stdErrors "errors"
	"io"
	"sync"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type Document struct {
	file    *fs.File
	scanner *scanner

	checkErr error

	len uint

	lenOnce   sync.Once
	checkOnce sync.Once

	allowTrailingNonSpaceCharacters bool
}

var _ jschema.Document = &Document{}

// New creates a JSON document with specified name and content.
func New(name string, content []byte, oo ...Option) jschema.Document {
	return FromFile(fs.NewFile(name, content), oo...)
}

// FromFile creates a JSON document from file.
func FromFile(f *fs.File, oo ...Option) jschema.Document {
	d := &Document{
		file: f,
	}

	for _, o := range oo {
		o(d)
	}

	d.rewind()

	return d
}

type Option func(s *Document)

func AllowTrailingNonSpaceCharacters() Option {
	return func(s *Document) {
		s.allowTrailingNonSpaceCharacters = true
	}
}

func (d *Document) NextLexeme() (lexeme.LexEvent, error) {
	return d.nextLexeme()
}

func (d *Document) Len() (uint, error) {
	var err error
	d.lenOnce.Do(func() {
		d.len, err = d.computeLen()
	})
	return d.len, err
}

func (d *Document) computeLen() (length uint, err error) {
	// Iterate through all lexemes until we reach the end
	// We should rewind here in case we call NextLexeme method.
	d.rewind()
	defer d.rewind()
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		rErr, ok := r.(error)
		if !ok {
			panic(r)
		}
		err = rErr
	}()

	return d.scanner.Length(), err
}

func (d *Document) Check() error {
	d.checkOnce.Do(func() {
		d.checkErr = d.check()
	})
	return d.checkErr
}

func (d *Document) check() error {
	// Iterate through all lexemes until we reach the end or get some error.
	// We should rewind here in case we call NextLexeme method.
	d.rewind()
	defer d.rewind()

	var jsonLexCounter uint
	for {
		_, err := d.nextLexeme()
		if err == nil {
			jsonLexCounter++
			continue
		}

		if stdErrors.Is(err, io.EOF) {
			err = nil

			if jsonLexCounter == 0 {
				err = errors.NewDocumentError(d.file, errors.ErrEmptyJson)
			}
		}
		return err
	}
}

func (d *Document) nextLexeme() (lex lexeme.LexEvent, err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		rErr, ok := r.(error)
		if !ok {
			panic(r)
		}
		err = rErr
	}()

	lex, ok := d.scanner.Next()
	if !ok {
		return lexeme.LexEvent{}, io.EOF
	}

	if lex.Type() == lexeme.EndTop {
		return lex, io.EOF
	}
	return lex, nil
}

// rewind rewinds document to the beginning.
func (d *Document) rewind() {
	d.scanner = newScanner(d.file)
	d.scanner.allowTrailingNonSpaceCharacters = d.allowTrailingNonSpaceCharacters
}
