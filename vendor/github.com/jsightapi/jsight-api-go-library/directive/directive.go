package directive

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/jerr"
)

// Directive represents all info about some Directive
type Directive struct {
	Annotation string

	includeTracer IncludeTracer

	// Keyword only for Responses (have multiple keywords), for others should match
	// type.
	Keyword       string
	parameters    map[string]string
	keywordCoords Coords
	BodyCoords    Coords
	Parent        *Directive
	Children      []*Directive

	type_ Enumeration

	// HasExplicitContext true if directive's context is opened explicitly with parentheses.
	HasExplicitContext bool
}

func New(e Enumeration, keywordCoords Coords) *Directive {
	return NewWithCallStack(e, keywordCoords, nopIncludeTracer{})
}

func NewWithCallStack(e Enumeration, keywordCoords Coords, includeTracer IncludeTracer) *Directive {
	return &Directive{
		type_:         e,
		parameters:    make(map[string]string),
		Keyword:       e.String(),
		keywordCoords: keywordCoords,
		includeTracer: includeTracer,
	}
}

func (d Directive) String() string {
	return d.Type().String()
}

func (d Directive) Equal(d2 Directive) bool {
	return d.keywordCoords.file == d2.keywordCoords.file &&
		d.keywordCoords.begin == d2.keywordCoords.begin
}

func (d Directive) Type() Enumeration {
	return d.type_
}

func (d Directive) HasAnyParameters() bool {
	return len(d.parameters) != 0
}

func (d Directive) Parameter(k string) string {
	if v, ok := d.parameters[k]; ok {
		return v
	}
	return ""
}

func (d *Directive) SetParameter(k string, v string) error {
	if _, ok := d.parameters[k]; ok {
		return fmt.Errorf("the %q parameter is already defined for the %q directive", k, d)
	}
	d.parameters[k] = v
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
		parameters:         d.parameters,
		keywordCoords:      d.keywordCoords,
		BodyCoords:         d.BodyCoords,
		includeTracer:      d.includeTracer,
	}
}

// IncludeTracer represent the directive's call stack.
type IncludeTracer interface {
	// AddIncludeTraceToError adds proper trace to error.
	AddIncludeTraceToError(je *jerr.JApiError)
}

type nopIncludeTracer struct{}

func (nopIncludeTracer) AddIncludeTraceToError(*jerr.JApiError) {}
