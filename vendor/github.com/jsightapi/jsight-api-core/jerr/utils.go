package jerr

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
)

// GetQuote return "" if the source sub-string cannot be determined
func GetQuote(content bytes.Bytes, position bytes.Index) string {
	return quote(content, position)
}

func quote(content bytes.Bytes, position bytes.Index) string {
	const maxLength = 200
	begin := content.BeginningOfLine(position)
	end := content.EndOfLine(position)
	if end-begin > maxLength {
		end = begin + maxLength - 3
		return content.Sub(begin, end).TrimSpacesFromLeft().String() + "..."
	}
	return content.Sub(begin, end).TrimSpacesFromLeft().String()
}
