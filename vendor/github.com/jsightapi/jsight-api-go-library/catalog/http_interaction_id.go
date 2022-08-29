package catalog

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
)

type HTTPInteractionID struct {
	protocol Protocol
	path     Path
	method   HTTPMethod
}

func (h HTTPInteractionID) Protocol() Protocol {
	return h.protocol
}

func (h HTTPInteractionID) Path() Path {
	return h.path
}

func (h HTTPInteractionID) String() string {
	return fmt.Sprintf("http %s %s", h.method.String(), h.path.String())
}

func (h HTTPInteractionID) MarshalText() ([]byte, error) {
	return []byte(h.String()), nil
}

func newHTTPInteractionID(d directive.Directive) (HTTPInteractionID, error) {
	h := HTTPInteractionID{
		protocol: HTTP,
	}

	path, err := d.Path()
	if err != nil {
		return h, err
	}

	de, err := d.HTTPMethod()
	if err != nil {
		return h, err
	}

	method, err := NewHTTPMethod(de)
	if err != nil {
		return h, err
	}

	h.path = Path(path)
	h.method = method

	return h, nil
}
