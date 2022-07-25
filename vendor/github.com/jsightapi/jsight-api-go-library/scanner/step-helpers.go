package scanner

func caseWhitespace(c byte) byte {
	if isWhitespace(c) {
		return c
	} else {
		return otherByte(c)
	}
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t'
}

func caseNewLine(c byte) byte {
	if IsNewLine(c) {
		return c
	} else {
		return otherByte(c)
	}
}

func IsNewLine(c byte) bool {
	return c == '\n' || c == '\r'
}

// byte is uint8: 0-255. 0 is for EOF.
func otherByte(b byte) byte {
	if b == 255 {
		return 254
	} else {
		return b + 1
	}
}

var (
	anyType   = []byte("any")
	emptyType = []byte("empty")
	regexType = []byte("regex")
)

func (s *Scanner) isDirectiveParameterHasTypeOrAnyOrEmpty() bool {
	for _, lex := range s.lastDirectiveParameters {
		v := lex.Value().Unquote().TrimSquareBrackets()
		switch {
		case v.Equals(anyType), v.Equals(emptyType), v.IsUserTypeName():
			return true
		}
	}
	return false
}

func (s *Scanner) isDirectiveParameterHasAnyOrEmpty() bool {
	for _, lex := range s.lastDirectiveParameters {
		v := lex.Value().Unquote().TrimSquareBrackets()
		switch {
		case v.Equals(anyType), v.Equals(emptyType):
			return false
		}
	}
	return true
}

func (s *Scanner) isDirectiveParameterHasRegexNotation() bool {
	for _, lex := range s.lastDirectiveParameters {
		v := lex.Value().Unquote()
		if v.Equals(regexType) {
			return true
		}
	}
	return false
}
