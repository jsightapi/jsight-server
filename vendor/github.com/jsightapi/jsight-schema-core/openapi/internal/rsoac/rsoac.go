package rsoac

import "github.com/jsightapi/jsight-schema-core/notations/regex"

// Regex schema to OpenAPi converter

type RSOAC struct {
	description *string
}

func New(rs *regex.RSchema) *RSOAC {
	panic("TODO regex.RSchema") // TODO func
}

func (o *RSOAC) SetDescription(s string) {
	o.description = &s
}

func (o RSOAC) MarshalJSON() (b []byte, err error) {
	return []byte("TODO"), nil // TODO method
}
