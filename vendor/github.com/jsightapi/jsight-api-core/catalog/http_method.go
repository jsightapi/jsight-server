package catalog

import (
	"errors"

	"github.com/jsightapi/jsight-api-core/directive"
	"github.com/jsightapi/jsight-api-core/jerr"
)

type HTTPMethod uint8

const (
	GET HTTPMethod = iota
	POST
	PUT
	PATCH
	DELETE
	OPTIONS
)

func NewHTTPMethod(de directive.Enumeration) (HTTPMethod, error) {
	switch de {
	case directive.Get:
		return GET, nil
	case directive.Post:
		return POST, nil
	case directive.Put:
		return PUT, nil
	case directive.Patch:
		return PATCH, nil
	case directive.Delete:
		return DELETE, nil
	default:
		return GET, errors.New(jerr.RuntimeFailure)
	}
}

func NewHTTPMethodFromString(s string) (HTTPMethod, error) {
	switch s {
	case directive.Get.String():
		return GET, nil
	case directive.Post.String():
		return POST, nil
	case directive.Put.String():
		return PUT, nil
	case directive.Patch.String():
		return PATCH, nil
	case directive.Delete.String():
		return DELETE, nil
	default:
		return GET, errors.New(jerr.RuntimeFailure)
	}
}

func (e HTTPMethod) String() string {
	switch e {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case PATCH:
		return "PATCH"
	case DELETE:
		return "DELETE"
	case OPTIONS:
		return "OPTIONS"
	default:
		panic("Unknown method")
	}
}

func (e HTTPMethod) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}
