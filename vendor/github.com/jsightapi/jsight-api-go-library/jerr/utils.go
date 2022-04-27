package jerr

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

func DetectNewLineSymbol(content bytes.Bytes) byte {
	newLineByte := byte('\n') // default new line
	var found bool
	var i = 0
	for i < len(content) {
		c := content[i]
		if c == '\n' || c == '\r' {
			newLineByte = c
			found = true
		} else if found { // first symbol after new line
			break
		}
		i++
	}
	return newLineByte
}

func PositionInLine(content bytes.Bytes, position bytes.Index, nl byte) bytes.Index {
	lb := LineBeginning(content, position, nl)
	return position - lb
}

// GetQuote return "" if cannot determine the source sub-string
func GetQuote(content bytes.Bytes, position bytes.Index, nl byte) string {
	begin := LineBeginning(content, position, nl)
	return quote(content, position, begin, nl)
}

func quote(content bytes.Bytes, position bytes.Index, lineBeginning bytes.Index, nl byte) string {
	end := LineEnd(content, position, nl)
	maxLength := bytes.Index(200)
	if end-lineBeginning > maxLength {
		end = lineBeginning + maxLength - 3
		return string(content[lineBeginning:end].TrimSpacesFromLeft()) + "..."
	}
	return string(content[lineBeginning:end].TrimSpacesFromLeft())
}

// LineNumber return 0 if cannot determine the line number, or 1+ if it can
func LineNumber(content bytes.Bytes, position bytes.Index, nl byte) bytes.Index {
	i := position
	max := bytes.Index(len(content) - 1)
	if i > max {
		i = max
	}
	var n uint
	for {
		c := content[i]
		if c == nl {
			if i != position {
				n++
			}
		}
		if i == 0 { // It is important because an unsigned value (i := 0; i--; i == [large positive number])
			break
		}
		i--
	}
	return bytes.Index(n + 1)
}

// LineBeginning Before calling this method, you must run the e.preparation()
func LineBeginning(content bytes.Bytes, position bytes.Index, nl byte) bytes.Index {
	i := position
	max := bytes.Index(len(content) - 1)
	if i > max {
		i = max
	}
	for {
		c := content[i]
		if c == nl {
			if i != position {
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

func LineEnd(content bytes.Bytes, position bytes.Index, nl byte) bytes.Index {
	i := position
	for i < bytes.Index(len(content)) {
		c := content[i]
		if c == nl {
			break
		}
		i++
	}
	if i > 0 {
		c := content[i-1]
		if (nl == '\n' && c == '\r') || (nl == '\r' && c == '\n') {
			i--
		}
	}
	return i
}
