package directive

import (
	"fmt"
)

// Directive represents all info about some Directive
type Directive struct {
	Annotation string

	// Keyword only for Responses (have multiple keywords), for others should match
	// type.
	Keyword       string
	parameter     map[string]string
	keywordCoords Coords
	BodyCoords    Coords
	Parent        *Directive
	Children      []*Directive

	type_ Enumeration

	// HasExplicitContext true if directive's context is opened explicitly with parentheses.
	HasExplicitContext bool
}

func (d Directive) String() string {
	return d.Type().String()
}

func (d Directive) Equal(d2 Directive) bool {
	return d.keywordCoords.f == d2.keywordCoords.f && d.keywordCoords.b == d2.keywordCoords.b
}

func (d Directive) Type() Enumeration {
	return d.type_
}

func (d Directive) HasAnyParameters() bool {
	return len(d.parameter) != 0
}

func (d Directive) Parameter(k string) string {
	if v, ok := d.parameter[k]; ok {
		return v
	}
	return ""
}

func (d *Directive) SetParameter(k string, v string) error {
	if _, ok := d.parameter[k]; ok {
		return fmt.Errorf("the %q parameter is already defined for the %q directive", k, d)
	}
	d.parameter[k] = v
	return nil
}

func (d *Directive) AppendChild(child *Directive) {
	if d.Children == nil {
		d.Children = make([]*Directive, 0, 10)
	}
	d.Children = append(d.Children, child)
}

func (d Directive) CopyWoParentAndChildren() Directive {
	return Directive{
		type_:              d.type_,
		Annotation:         d.Annotation,
		Keyword:            d.Keyword,
		HasExplicitContext: d.HasExplicitContext,
		parameter:          d.parameter,
		keywordCoords:      d.keywordCoords,
		BodyCoords:         d.BodyCoords,
		// Children:           nil,
		// Parent:             nil,
	}
}

func New(e Enumeration, keywordCoords Coords) *Directive {
	return &Directive{
		type_:         e,
		parameter:     make(map[string]string),
		Keyword:       e.String(),
		keywordCoords: keywordCoords,
	}
}
