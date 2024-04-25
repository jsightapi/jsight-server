package jsoac

import (
	"encoding/json"
)

type Number struct {
	value []byte
}

var _ json.Marshaler = Number{}
var _ json.Marshaler = &Number{}

// newNumber creates a Number for: integer, float, double values
func newNumber(n string) *Number {
	return &Number{value: []byte(n)}
}

func (n Number) MarshalJSON() (b []byte, err error) {
	return n.value, nil
}
