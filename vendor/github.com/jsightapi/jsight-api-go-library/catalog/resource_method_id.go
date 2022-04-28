package catalog

import (
	"github.com/jsightapi/jsight-api-go-library/directive"
)

type ResourceMethodId struct {
	path   Path
	method Method
}

func newResourceMethodId(d directive.Directive) (ResourceMethodId, error) {
	rk := ResourceMethodId{}

	path, err := d.Path()
	if err != nil {
		return rk, err
	}

	de, err := d.Method()
	if err != nil {
		return rk, err
	}

	method, err := NewMethod(de)
	if err != nil {
		return rk, err
	}

	rk.path = Path(path)
	rk.method = method

	return rk, nil
}

func (r ResourceMethodId) String() string {
	return r.method.String() + " " + r.path.String()
}

func (r ResourceMethodId) MarshalText() ([]byte, error) {
	return []byte(r.String()), nil
}
