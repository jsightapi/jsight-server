package core

import (
	"bytes"
)

func description(b []byte) []byte {
	b = descriptionTrimBrackets(b)
	b = bytes.TrimLeft(b, "\r\n")
	b = bytes.TrimRight(b, "\r\n\t ")

	lines := splitLines(b)
	prefix := longestWhitespacePrefix(lines)
	for i := 0; i < len(lines); i++ {
		lines[i] = bytes.TrimPrefix(lines[i], prefix)
	}
	return bytes.Join(lines, []byte{'\n'})
}

func descriptionTrimBrackets(b []byte) []byte {
	bb := bytes.TrimSpace(b)
	if len(bb) >= 2 && bb[0] == '(' && bb[len(bb)-1] == ')' {
		return bb[1 : len(bb)-1]
	}
	return b
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

func splitLines(b []byte) [][]byte {
	b = bytes.ReplaceAll(b, []byte{'\n', '\r'}, []byte{'\n'}) // Windows
	b = bytes.ReplaceAll(b, []byte{'\r'}, []byte{'\n'})       // Macintosh (old)
	return bytes.Split(b, []byte{'\n'})
}
