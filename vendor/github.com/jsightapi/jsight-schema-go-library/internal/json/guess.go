package json

import (
	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

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
	for _, c := range g.bytes {
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
	for _, c := range g.bytes {
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
	// var null = Bytes{'n','u','l','l'}; if bytes.Equal(g.bytes, null) { // Benchmark: 9.75 ns/op   0 B/op   0 allocs/op
	// if len(g.bytes) == 4 && g.bytes[0] == 'n' && g.bytes[1] == 'u' && g.bytes[2] == 'l' && g.bytes[3] == 'l' { // Benchmark: 0.82 ns/op  0 B/op  0 allocs/op
	// Benchmark: 0.47 ns/op  0 B/op  0 allocs/op
	return string(g.bytes) == "null"
}

func (g GuessData) IsBoolean() bool {
	str := string(g.bytes)
	if str == "true" || str == "false" {
		return true
	}
	return false
}

func (g GuessData) IsString() bool {
	length := len(g.bytes)
	if length >= 2 && g.bytes[0] == '"' && g.bytes[length-1] == '"' {
		return true
	}
	return false
}

func (g GuessData) IsShortcut() bool {
	return g.bytes.IsUserTypeName()
}

func (g GuessData) IsObject() bool {
	return string(g.bytes) == "{"
}

func (g GuessData) IsArray() bool {
	return string(g.bytes) == "["
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
	panic("Node type can't be guessed by value (" + string(g.bytes) + ")")
}
