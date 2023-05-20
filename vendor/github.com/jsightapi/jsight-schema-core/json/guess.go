package json

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
)

const nullStr = "null"

type GuessData struct {
	number *Number
	bytes  bytes.Bytes
}

func Guess(b bytes.Bytes) *GuessData {
	g := GuessData{bytes: b}
	// the number is not built immediately, but only if necessary (for optimization)
	return &g
}

func (g *GuessData) Number() (*Number, error) {
	if g.number == nil {
		n, err := NewNumber(g.bytes)
		if err != nil {
			return nil, err
		}
		g.number = n
	}
	return g.number, nil
}

func (g *GuessData) IsInteger() bool {
	dot := false
	exp := false
	for _, c := range g.bytes.Data() {
		switch c {
		case '.':
			dot = true
		case 'e', 'E':
			exp = true
		}
	}
	if dot && !exp {
		return false
	}

	n, err := g.Number()
	if err != nil {
		return false
	}

	if n.LengthOfFractionalPart() != 0 {
		return false
	}

	return true
}

func (g *GuessData) IsFloat() bool {
	dot := false
	exp := false
	for _, c := range g.bytes.Data() {
		switch c {
		case '.':
			dot = true
		case 'e', 'E':
			exp = true
		}
	}
	if dot && !exp {
		return true
	}

	n, err := g.Number()
	if err != nil {
		return false
	}

	if n.LengthOfFractionalPart() != 0 {
		return true
	}

	return false
}

func (g GuessData) IsNull() bool {
	// Benchmark: 9.75 ns/op   0 B/op   0 allocs/op
	// var null = Bytes{'n','u','l','l'}; if bytes.Equal(g.bytes, null) {

	// Benchmark: 0.82 ns/op  0 B/op  0 allocs/op
	// if len(g.bytes) == 4 && g.bytes[0] == 'n' && g.bytes[1] == 'u' && g.bytes[2] == 'l' && g.bytes[3] == 'l' {

	// Benchmark: 0.47 ns/op  0 B/op  0 allocs/op
	return g.bytes.String() == nullStr
}

func (g GuessData) IsBoolean() bool {
	str := g.bytes.String()
	if str == "true" || str == "false" {
		return true
	}
	return false
}

func (g GuessData) IsString() bool {
	if g.bytes.Len() >= 2 && g.bytes.FirstByte() == '"' && g.bytes.LastByte() == '"' {
		return true
	}
	return false
}

func (g GuessData) IsShortcut() bool {
	return g.bytes.IsUserTypeName()
}

func (g GuessData) IsObject() bool {
	return g.bytes.String() == "{"
}

func (g GuessData) IsArray() bool {
	return g.bytes.String() == "["
}

func (g GuessData) JsonType() Type {
	if g.IsObject() {
		return TypeObject
	} else if g.IsArray() {
		return TypeArray
	}
	return g.LiteralJsonType()
}

func (g GuessData) LiteralJsonType() Type {
	switch {
	case g.IsString():
		return TypeString
	case g.IsBoolean():
		return TypeBoolean
	case g.IsNull():
		return TypeNull
	case g.IsInteger():
		return TypeInteger
	case g.IsFloat():
		return TypeFloat
	case g.IsShortcut():
		return TypeMixed
	}
	panic(errs.ErrNodeTypeCantBeGuessed.F(g.bytes.String()))
}
