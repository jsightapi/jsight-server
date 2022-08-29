package bytes

import (
	"bytes"
	"errors"
	"fmt"
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
	if len(b) >= 2 && b[0] == '"' && b[len(b)-1] == '"' {
		return unquoteBytes(b)
	}
	return b
}

func (b Bytes) TrimSquareBrackets() Bytes {
	lastCharIndex := len(b) - 1
	if lastCharIndex > 0 && b[0] == '[' && b[lastCharIndex] == ']' {
		return b[1:lastCharIndex]
	}
	return b
}

func (b Bytes) TrimSpaces() Bytes {
	blen := len(b)

	left := 0
	right := blen - 1

	for ; left < blen && IsBlank(b[left]); left++ {
	}

	if left >= blen {
		return Bytes{}
	}

	for ; right > 0 && IsBlank(b[right]); right-- {
	}

	return b[left : right+1]
}

func (b Bytes) TrimSpacesFromLeft() Bytes {
	for i, c := range b {
		if !IsBlank(c) {
			return b[i:]
		}
	}
	return b
}

func (b Bytes) CountSpacesFromLeft() int {
	for i, c := range b {
		if !IsBlank(c) {
			return i
		}
	}
	return 0
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
	if len(b) == 0 {
		return 0, errors.New("not enough data in ParseUint")
	}

	var u uint
	for _, c := range b {
		if !IsDigit(c) {
			return 0, fmt.Errorf("invalid byte (%s) found in ParseUint (%s)", string(c), b)
		}
		u = u*10 + uint(c-'0')
	}
	return u, nil
}

func (b Bytes) ParseInt() (int, error) {
	var (
		negative bool
		u        uint
		err      error
	)

	if b[0] == '-' {
		negative = true
		u, err = b[1:].ParseUint()
	} else {
		u, err = b.ParseUint()
	}

	if err != nil {
		return 0, err
	}

	if u > math.MaxInt {
		return 0, errors.New("too much data for int")
	}

	i := int(u)
	if negative {
		return -i, nil
	}
	return i, nil
}

func (b Bytes) IsUserTypeName() bool {
	if len(b) < 2 || b[0] != '@' {
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

func (b Bytes) Normalize() Bytes {
	return normalizeBytes(b)
}
