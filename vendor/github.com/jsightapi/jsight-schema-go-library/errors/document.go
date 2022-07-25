package errors

import (
	"fmt"
	"strings"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

// DocumentError contains methods for forming a detailed description of the error
// for a person.
// The resulting message will contain the filename, line number, and where the error
// occurred.
type DocumentError struct {
	// A file containing jSchema or JSON data.
	file              *fs.File
	message           string
	incorrectUserType string
	code              ErrorCode

	// index of the byte in which the error was found.
	index bytes.Index

	// A length of file content.
	length bytes.Index

	// hasIndex true if the value for Index have been defined.
	hasIndex bool

	// prepared is true when preliminary calculations are made, the results of
	// which are used in some methods.
	prepared bool

	// nl represent new line symbol.
	nl byte
}

var (
	_ Error = DocumentError{}
	_ error = DocumentError{}
)

func NewDocumentError(file *fs.File, err Err) DocumentError {
	return DocumentError{
		code:    err.Code(),
		message: err.Error(),
		file:    file,
	}
}

func (e DocumentError) Code() ErrorCode {
	return e.code
}

func (e DocumentError) ErrCode() int {
	return int(e.code)
}

func (e DocumentError) Filename() string {
	if e.file == nil {
		return ""
	}
	return e.file.Name()
}

func (e DocumentError) Message() string {
	return e.message
}

func (e DocumentError) Position() uint {
	return uint(e.index)
}

func (e DocumentError) Index() bytes.Index {
	return e.index
}

func (e *DocumentError) SetIndex(index bytes.Index) {
	e.index = index
	e.hasIndex = true
}

func (e DocumentError) IncorrectUserType() string {
	return e.incorrectUserType
}

func (e *DocumentError) SetIncorrectUserType(s string) {
	e.incorrectUserType = s
}

func (e *DocumentError) SetFile(file *fs.File) {
	e.file = file
}

func (e *DocumentError) SetMessage(message string) {
	e.message = message
}

// The method performs preparatory calculations, the results of which are used in other methods.
func (e *DocumentError) preparation() {
	if e.prepared {
		return
	}

	if e.file == nil {
		panic("The file is not specified")
	}
	e.length = bytes.Index(len(e.file.Content()))
	e.detectNewLineSymbol()
	e.prepared = true
}

func (e *DocumentError) detectNewLineSymbol() {
	content := e.file.Content()
	e.nl = '\n' // default new line
	var found bool
	for _, c := range content {
		if bytes.IsNewLine(c) {
			e.nl = c
			found = true
		} else if found { // first symbol after new line
			break
		}
	}
}

// Before calling this method, you must run the e.preparation()
func (e DocumentError) lineBeginning() bytes.Index {
	content := e.file.Content()
	i := e.index
	for {
		c := content[i]
		if c == e.nl {
			if i != e.index {
				i++ // step forward from new line
				break
			}
		}
		if i == 0 { // It is important because an unsigned value (i := 0; i--; i == [large positive number])
			break
		}
		i--
	}
	return i
}

// Before calling this method, you must run the e.preparation()
func (e DocumentError) lineEnd() bytes.Index {
	content := e.file.Content()
	i := e.index
	for i < e.length {
		c := content[i]
		if c == e.nl {
			break
		}
		i++
	}
	if i > 0 {
		c := content[i-1]
		if (e.nl == '\n' && c == '\r') || (e.nl == '\r' && c == '\n') {
			i--
		}
	}
	return i
}

// Line returns 0, if cannot determine the line number, or 1+ if it can.
func (e *DocumentError) Line() uint {
	if e.file == nil || len(e.file.Content()) == 0 {
		return 0
	}

	e.preparation()

	content := e.file.Content()
	i := e.index
	var n uint

	for {
		c := content[i]
		if c == e.nl {
			if i != e.index {
				n++
			}
		}
		if i == 0 { // It is important because an unsigned value (i := 0; i--; i == [large positive number])
			break
		}
		i--
	}

	return n + 1
}

// SourceSubString returns empty string, if cannot determine the source sub-string.
func (e *DocumentError) SourceSubString() string {
	if e.file == nil || len(e.file.Content()) == 0 {
		return ""
	}

	e.preparation()

	content := e.file.Content()
	begin := e.lineBeginning()
	end := e.lineEnd()
	maxLength := bytes.Index(200)

	if end-begin > maxLength {
		end = begin + maxLength - 3
		return string(content[begin:end].TrimSpacesFromLeft()) + "..."
	}

	return string(content[begin:end].TrimSpacesFromLeft())
}

func (e *DocumentError) pointerToTheErrorCharacter() string {
	e.preparation()

	content := e.file.Content()
	begin := e.lineBeginning()
	spaces := content[begin:].CountSpacesFromLeft()

	i := int(e.index) - int(begin) - spaces
	return strings.Repeat("-", i) + "^"
}

func (e DocumentError) Error() string {
	return e.String()
}

func (e *DocumentError) String() string {
	var prefix string
	if e.code == ErrGeneric {
		prefix = "ERROR"
	} else {
		prefix = "ERROR (code " + e.code.Itoa() + ")"
	}
	if e.file != nil {
		filename := e.file.Name()
		if e.hasIndex {
			return fmt.Sprintf(`%s: %s
	in line %d on file %s
	> %s
	--%s`, prefix, e.message, e.Line(), filename, e.SourceSubString(), e.pointerToTheErrorCharacter())
		} else if filename != "" {
			return fmt.Sprintf("%s: %s\n\tin file %s", prefix, e.message, filename)
		}
	}
	return fmt.Sprintf("%s: %s", prefix, e.message)
}
