package jerr

import (
	"strconv"
	"strings"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

type JApiError struct {
	Msg          string
	wrapped      *JApiError
	includeTrace []stackTraceItem

	Location
}

type stackTraceItem struct {
	path   string
	atLine bytes.Index
}

var _ error = &JApiError{}

func NewJApiError(msg string, f *fs.File, i bytes.Index) *JApiError {
	return &JApiError{
		Location: NewLocation(f, i),
		Msg:      msg,
	}
}

func (e *JApiError) OccurredInFile(f *fs.File, atByte bytes.Index) {
	e.includeTrace = append(e.includeTrace, stackTraceItem{
		path:   f.Name(),
		atLine: NewLocation(f, atByte).line,
	})
}

func (e *JApiError) HasStackTrace() bool {
	return e != nil && len(e.includeTrace) > 0
}

func (e *JApiError) Error() string {
	if len(e.includeTrace) == 0 {
		return e.Msg
	}

	return e.errorWithStackTrace()
}

func (e *JApiError) errorWithStackTrace() string {
	buf := strings.Builder{}

	buf.WriteString(e.Msg)

	writeStackTraceLine(&buf, e.file.Name(), e.line)
	for _, i := range e.includeTrace {
		writeStackTraceLine(&buf, i.path, i.atLine)
	}

	return buf.String()
}

func writeStackTraceLine(buf *strings.Builder, p string, atLine bytes.Index) {
	buf.WriteRune('\n')
	buf.WriteString(p)
	buf.WriteRune(':')
	buf.WriteString(strconv.FormatUint(uint64(atLine), 10))
}
