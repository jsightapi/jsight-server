package bytes

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
