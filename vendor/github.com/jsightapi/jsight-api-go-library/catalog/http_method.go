package catalog

import (
	"errors"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

type HttpMethod uint8

const (
	GET HttpMethod = iota
	POST
	PUT
	PATCH
	DELETE
	OPTIONS
)

func NewHttpMethod(de directive.Enumeration) (HttpMethod, error) {
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
		return GET, errors.New(jerr.IsNotHTTPRequestMethod)
	}
}

func (e HttpMethod) String() string {
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

func (e HttpMethod) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}
