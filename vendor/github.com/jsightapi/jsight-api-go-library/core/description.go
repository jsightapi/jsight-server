package core

import (
	"bytes"
	"errors"

	"github.com/jsightapi/jsight-api-go-library/jerr"
	"github.com/jsightapi/jsight-api-go-library/scanner"
)

func description(b []byte) ([]byte, error) {
	b = bytes.ReplaceAll(b, []byte{'\r', '\n'}, []byte{'\n'}) // Windows
	b = bytes.ReplaceAll(b, []byte{'\r'}, []byte{'\n'})       // Macintosh (old)

	b, err := descriptionRemoveParentheses(b)
	if err != nil {
		return b, err
	}

	b = bytes.TrimLeft(b, "\r\n")
	b = bytes.TrimRight(b, "\r\n\t ")

	lines := bytes.Split(b, []byte{'\n'})

	prefix := longestWhitespacePrefix(lines)
	for i := 0; i < len(lines); i++ {
		lines[i] = bytes.TrimPrefix(lines[i], prefix)
	}

	return bytes.Join(lines, []byte{'\n'}), nil
}

// descriptionRemoveParentheses - removes parentheses and their accompanying whitespace and new line characters. Checks
// the text inside the parentheses for compliance with the rules:
//
// - The opening parenthesis ( must be placed on a new line immediately after the line containing the KEYWORD of the
// DIRECTIVE, or through any number of empty lines, since empty lines are ignored by the parser. This parenthesis
// declares the beginning of the BODY of the DIRECTIVE. Apart from the opening parenthesis, there should be nothing
// else on this line.
//
// - The closing parenthesis ) must be placed on a separate line immediately after the contents of the BODY of the
// DIRECTIVE or through any number of empty lines, since empty lines are ignored by the parser. The closing parenthesis
// declares the end of the BODY of the DIRECTIVE. Except the closing parenthesis, there should be nothing else on this
// line.
func descriptionRemoveParentheses(b []byte) ([]byte, error) {
	bb := bytes.TrimSpace(b) // trim whitespaces and new lines (outside parentheses)
	if len(bb) >= 2 && bb[0] == '(' && bb[len(bb)-1] == ')' {
		bb = bb[1 : len(bb)-1]     // trim parentheses
		bb = bytes.Trim(bb, " \t") // trim whitespaces (inside parentheses)
		if len(bb) == 0 || !scanner.IsNewLine(bb[0]) || !scanner.IsNewLine(bb[len(bb)-1]) {
			return bb, errors.New(jerr.ApartFromTheOpeningParenthesis)
		}
		return bytes.Trim(bb, "\r\n"), nil
	}
	return b, nil
}

func longestWhitespacePrefix(bb [][]byte) []byte {
	empty := make([]byte, 0)

	if len(bb) == 0 {
		return empty
	}

	prefix := empty
	for i := 0; i < len(bb[0]); i++ {
		if bb[0][i] != '\t' && bb[0][i] != ' ' || i == len(bb[0])-1 {
			prefix = bb[0][:i]
			break
		}
	}

	if len(prefix) == 0 {
		return empty
	}

	for i := 1; i < len(bb); i++ {
		if len(bb[i]) != 0 {
			for !bytes.HasPrefix(bb[i], prefix) {
				prefix = prefix[:len(prefix)-1]
				if len(prefix) == 0 {
					return empty
				}
			}
		}
	}

	return prefix
}
