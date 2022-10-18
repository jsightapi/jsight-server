package bytes

import "strconv"

// IsBlank returns true if provided byte is space or a new line.
func IsBlank(c byte) bool {
	return IsSpace(c) || IsNewLine(c)
}

// IsSpace returns true is provided byte is space.
func IsSpace(c byte) bool {
	return c == ' ' || c == '\t'
}

// IsNewLine returns true if provided byte is a new line.
func IsNewLine(c byte) bool {
	return c == '\n' || c == '\r'
}

// IsDigit returns true if provided byte is a digit.
func IsDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// IsHexDigit returns true if provided byte is a hex digit.
func IsHexDigit(c byte) bool {
	return IsDigit(c) || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// IsValidUserTypeNameByte returns true if specified rune can be a part of user
// type name.
// Important: `@` isn't valid here 'cause schema name should start with `@` but
// it didn't allow to use that symbol in the name.
func IsValidUserTypeNameByte(c byte) bool {
	return c == '-' || c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || IsDigit(c)
}

// QuoteChar formats char as a quoted character literal.
func QuoteChar(c byte) string {
	// special cases - different from quoted strings
	if c == '\'' {
		return `'\''`
	}
	if c == '"' {
		return `'"'`
	}
	// use quoted string with different quotation marks
	s := strconv.Quote(string(c))
	return "'" + s[1:len(s)-1] + "'"
}
