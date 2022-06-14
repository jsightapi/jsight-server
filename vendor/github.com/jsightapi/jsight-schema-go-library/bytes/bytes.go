package bytes

import (
	"bytes"
	"errors"
	"math"
)

type Index uint
type Bytes []byte

func (b Bytes) Equals(bb Bytes) bool {
	return bytes.Equal(b, bb)
}

func (b Bytes) Slice(begin, end Index) Bytes {
	return b[begin : end+1]
}

func (b Bytes) Unquote() Bytes {
	if len(b) >= 2 {
		lestCharIndex := len(b) - 1
		if b[0] == '"' && b[lestCharIndex] == '"' {
			return b[1:lestCharIndex]
		}
	}
	return b
}

func (b Bytes) TrimSquareBrackets() Bytes {
	if len(b) >= 2 {
		lestCharIndex := len(b) - 1
		if b[0] == '[' && b[lestCharIndex] == ']' {
			return b[1:lestCharIndex]
		}
	}
	return b
}

func (b Bytes) TrimSpaces() Bytes {
	blen := len(b)

	left := 0
	right := blen - 1

	for ; left < blen && isSpace(b[left]); left++ {
	}

	if left >= blen {
		return Bytes{}
	}

	for ; right > 0 && isSpace(b[right]); right-- {
	}

	return b[left : right+1]
}

func (b Bytes) TrimSpacesFromLeft() Bytes {
	for i, c := range b {
		if !isSpace(c) {
			return b[i:]
		}
	}
	return b
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func (b Bytes) CountSpacesFromLeft() int {
	var i int
	var c byte
	for i, c = range b {
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
	}
	return i
}

// OneOf checks current bytes sequence equal to at least one of specified strings.
func (b Bytes) OneOf(ss ...string) bool {
	// It's the fastest solution.
	for _, s := range ss {
		if string(b) == s {
			return true
		}
	}
	return false
}

// func (b Bytes) TrimZeroFromRight() Bytes {
// 	length := len(b)
// 	if length > 0 {
// 		cut := -1
// 		for i := length - 1; i >= 0; i-- {
// 			c := b[i]
// 			if c == '0' {
// 				cut = i
// 			} else {
// 				break
// 			}
// 		}
// 		if cut != -1 {
// 			return b[:cut]
// 		}
// 	}
// 	return b
// }

func (b Bytes) ParseBool() (bool, error) {
	switch string(b) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, errors.New("invalid bool value")
}

func (b Bytes) ParseUint() (uint, error) {
	var u uint
	if len(b) == 0 {
		return 0, errors.New("not enough data in ParseUint")
	}
	for _, c := range b {
		if '0' <= c && c <= '9' { //nolint:revive // early-return: if c {...} else {... return } can be simplified to if !c { ... return } ...
			c -= '0'
			u = u*10 + uint(c)
		} else {
			return 0, errors.New("invalid byte (" + string(c) + ") found in ParseUint (" + string(b) + ")")
		}
	}
	return u, nil
}

func (b Bytes) ParseInt() (int, error) {
	var negative bool // = false
	var u uint
	var err error

	if b[0] == '-' {
		negative = true
		u, err = b[1:].ParseUint()
	} else {
		u, err = b.ParseUint()
	}

	if err != nil {
		return 0, err
	}

	if u > math.MaxUint/2 {
		return 0, errors.New("too much data for int")
	}

	i := int(u)
	if negative {
		return -i, nil
	}
	return i, nil
}

func (b Bytes) IsUserTypeName() bool {
	if len(b) < 2 {
		return false
	}
	if b[0] != '@' {
		return false
	}
	for _, c := range b[1:] {
		if !IsValidUserTypeNameByte(c) {
			return false
		}
	}
	return true
}

func (b Bytes) String() string {
	return string(b)
}

func (b Bytes) Len() int {
	return len(b)
}

func (b Bytes) LineFrom(start Index) (Bytes, error) {
	l := Index(len(b))
	if start > l {
		return b, errors.New("can't get a line from a slice")
	}
	for i := start; i < l; i++ {
		if b[i] == '\n' {
			return b[start:i], nil
		}
	}
	return b[start:], nil
}

// IsValidUserTypeNameByte returns true if specified rune can be a part of user
// type name.
// Important: `@` isn't valid here 'cause schema name should start with `@` but
// it didn't allow to use that symbol in the name.
func IsValidUserTypeNameByte(c byte) bool {
	return c == '-' || c == '_' || ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}
