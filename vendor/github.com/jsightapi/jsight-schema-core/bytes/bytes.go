package bytes

import (
	"bytes"
	"math"
	"unicode/utf8"

	"github.com/jsightapi/jsight-schema-core/errs"
)

type Bytes struct {
	data []byte
	nl   byte
}

// ByteKeeper all allowed types for specifying Byte constructor
type ByteKeeper interface {
	string | []byte | Bytes
}

func MakeBytes(size int) Bytes {
	return Bytes{
		data: make([]byte, 0, size),
	}
}

func NewBytes[T ByteKeeper](data T) Bytes {
	switch bb := any(data).(type) {
	case nil:
		return Bytes{data: nil}
	case []byte:
		return Bytes{data: bb}
	case string:
		return Bytes{data: []byte(bb)}
	case Bytes:
		return bb
	}
	// This might happen only when we extend `ByteKeeper` interface and forget
	// to add new case to the type switch above this point.
	panic(errs.ErrRuntimeFailure.F())
}

func (b Bytes) IsNil() bool {
	return b.data == nil
}

func (b *Bytes) Append(c byte) {
	b.data = append(b.data, c)
}

func (b Bytes) Byte(i any) byte {
	return b.data[Int(i)]
}

func (b Bytes) FirstByte() byte {
	return b.data[0]
}

func (b Bytes) LastByte() byte {
	return b.data[b.Len()-1]
}

func (b Bytes) Data() []byte {
	return b.data
}

func (b Bytes) DecodeRune() rune {
	r, _ := utf8.DecodeRune(b.data)
	return r
}

func (b Bytes) ToLower() Bytes {
	return NewBytes(bytes.ToLower(b.data))
}

func (b Bytes) Sub(low, high any) Bytes {
	l := Int(low)
	h := Int(high)
	return NewBytes(b.data[l:h])
}

func (b Bytes) SubLow(low any) Bytes {
	i := Int(low)
	return NewBytes(b.data[i:])
}

func (b Bytes) SubHigh(high any) Bytes {
	i := Int(high)
	return NewBytes(b.data[:i])
}

func (b Bytes) SubToEndOfLine(start Index) (Bytes, error) {
	if start > b.LenIndex() {
		return b, errs.ErrRuntimeFailure.F()
	}

	bb := b.data[start:]

	for i, c := range bb {
		if c == '\n' || c == '\r' {
			bb = bb[:i]
			break
		}
	}

	return NewBytes(bb), nil
}

func (b Bytes) Equals(bb Bytes) bool {
	return bytes.Equal(b.data, bb.data)
}

// InQuotes the function is only needed in order not to modify the library function unquoteBytes()
func (b Bytes) InQuotes() bool {
	return len(b.data) >= 2 && b.data[0] == '"' && b.data[len(b.data)-1] == '"'
}

func (b Bytes) Unquote() Bytes {
	if b.InQuotes() {
		bb, ok := unquoteBytes(b.data)
		if !ok {
			return b // Can this happen?
		}
		return NewBytes(bb)
	}
	return b
}

func (b Bytes) TrimSquareBrackets() Bytes {
	lastCharIndex := len(b.data) - 1
	if lastCharIndex > 0 && b.data[0] == '[' && b.data[lastCharIndex] == ']' {
		return NewBytes(b.data[1:lastCharIndex])
	}
	return b
}

func (b Bytes) TrimSpaces() Bytes {
	left := 0
	right := b.Len() - 1

	for ; left < b.Len() && IsBlank(b.data[left]); left++ {
	}

	if left >= b.Len() {
		return Bytes{}
	}

	for ; right > 0 && IsBlank(b.data[right]); right-- {
	}

	return NewBytes(b.data[left : right+1])
}

func (b Bytes) TrimSpacesFromLeft() Bytes {
	for i, c := range b.data {
		if !IsBlank(c) {
			return NewBytes(b.data[i:])
		}
	}
	return b
}

func (b Bytes) CountSpacesFromLeft() int {
	for i, c := range b.data {
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
		if string(b.data) == s {
			return true
		}
	}
	return false
}

func (b Bytes) ParseBool() (bool, error) {
	switch string(b.data) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, errs.ErrInvalidBoolValue.F()
}

func (b Bytes) ParseUint() (uint, error) {
	if len(b.data) == 0 {
		return 0, errs.ErrNotEnoughDataInParseUint.F()
	}

	var u uint
	for _, c := range b.data {
		if !IsDigit(c) {
			return 0, errs.ErrInvalidByteInParseUint.F(string(c), b)
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

	if b.data[0] == '-' {
		negative = true
		u, err = NewBytes(b.data[1:]).ParseUint()
	} else {
		u, err = b.ParseUint()
	}

	if err != nil {
		return 0, err
	}

	if u > math.MaxInt {
		return 0, errs.ErrTooMuchDataForInt.F()
	}

	i := int(u)
	if negative {
		return -i, nil
	}
	return i, nil
}

func (b Bytes) IsUserTypeName() bool {
	if len(b.data) < 2 || b.data[0] != '@' {
		return false
	}
	for _, c := range b.data[1:] {
		if !IsValidUserTypeNameByte(c) {
			return false
		}
	}
	return true
}

func (b Bytes) String() string {
	return string(b.data)
}

func (b Bytes) Len() int {
	return len(b.data)
}

func (b Bytes) LenIndex() Index {
	return Index(b.Len())
}

// LineAndColumn calculate the line and column numbers by byte index in the content
// return 0 if not found
func (b Bytes) LineAndColumn(index Index) (line, column Index) {
	if b.Len() == 0 || Index(b.Len()) <= index {
		return 0, 0
	}

	nl := b.NewLineSymbol()

	for _, c := range b.data[:index] {
		if c == nl {
			line++
			column = 0
		} else {
			column++
		}
	}

	line++
	column++

	return line, column
}

func (b Bytes) BeginningOfLine(index Index) Index {
	nl := b.NewLineSymbol()
	i := index
	max := b.LenIndex() - 1
	if i > max {
		i = max
	}
	for {
		c := b.data[i]
		if c == nl {
			if i != index {
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

func (b Bytes) EndOfLine(index Index) Index {
	nl := b.NewLineSymbol()
	i := index
	for i < b.LenIndex() {
		c := b.data[i]
		if c == nl {
			break
		}
		i++
	}
	if i > 0 {
		c := b.data[i-1]
		if (nl == '\n' && c == '\r') || (nl == '\r' && c == '\n') {
			i--
		}
	}
	return i
}

func (b Bytes) NewLineSymbol() byte {
	switch b.nl {
	case '\n', '\r':
		return b.nl
	default:
		b.nl = '\n' // default
		var found bool
		for _, c := range b.data {
			if c == '\n' || c == '\r' {
				b.nl = c
				found = true
			} else if found { // first symbol after new line
				break
			}
		}
	}
	return b.nl
}
