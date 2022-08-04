package directive

import (
	"errors"
	"strings"
	"sync"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
)

type Enumeration int

const (
	Jsight Enumeration = iota
	Info
	Title
	Version
	Description
	Server
	BaseUrl
	Url
	Get
	Post
	Put
	Patch
	Delete
	Body
	Request
	HTTPResponseCode
	Path
	Headers
	Query
	Type
	Enum
	Macro
	Paste

	// Include directive.
	// This directive will be processed before we build a directive tree, so you
	// won't see it.
	Include

	Protocol
	Method
	Params
	Result
	TAG
	Tags
)

var (
	ss = []string{
		"JSIGHT",
		"INFO",
		"Title",
		"Version",
		"Description",
		"SERVER",
		"BaseUrl",
		"URL",
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"Body",
		"Request",
		"HTTP-response-code",
		"Path",
		"Headers",
		"Query",
		"TYPE",
		"ENUM",
		"MACRO",
		"PASTE",
		"INCLUDE",
		"Protocol",
		"Method",
		"Params",
		"Result",
		"TAG",
		"Tags",
	}
	eeOnce sync.Once
	ee     map[string]Enumeration
)

func (de Enumeration) String() string {
	return ss[de]
}

func NewDirectiveType(s string) (Enumeration, error) {
	eeOnce.Do(func() {
		ee = make(map[string]Enumeration)
		for i := 0; i < len(ss); i++ {
			if Enumeration(i) != HTTPResponseCode {
				ee[ss[i]] = Enumeration(i)
			}
		}
	})

	if v, ok := ee[s]; ok {
		return v, nil
	}

	if IsHTTPResponseCode(s) {
		return HTTPResponseCode, nil
	}

	return Jsight, errors.New("unknown directive type")
}

func (de Enumeration) IsHTTPRequestMethod() bool {
	switch de { //nolint:exhaustive // False-positive.
	case Get, Post, Put, Patch, Delete:
		return true
	}
	return false
}

func (de Enumeration) IsAllowedForRootContext() bool {
	switch de { //nolint:exhaustive // False-positive.
	case Jsight, Info, Server, Url, Get, Post, Put, Patch, Delete, Type, Enum,
		Macro, Paste, TAG:
		return true
	}
	return false
}

func (de Enumeration) IsAllowedForDirectiveContext(child Enumeration) bool {
	s, ok := directiveAllowedToDirectiveContext[de]
	if !ok {
		return false
	}

	_, ok = s[child]
	return ok
}

// directiveAllowedToDirectiveContext a map between directive type and directive
// types which can be placed into this directive context.
var directiveAllowedToDirectiveContext = map[Enumeration]map[Enumeration]struct{}{
	Url:              createEnumerationSet(Get, Post, Put, Patch, Delete, Path, Paste, Protocol, Method, Tags),
	Get:              createEnumerationSet(Description, Request, HTTPResponseCode, Path, Query, Paste, Tags),
	Post:             createEnumerationSet(Description, Request, HTTPResponseCode, Path, Query, Paste, Tags),
	Put:              createEnumerationSet(Description, Request, HTTPResponseCode, Path, Query, Paste, Tags),
	Patch:            createEnumerationSet(Description, Request, HTTPResponseCode, Path, Query, Paste, Tags),
	Delete:           createEnumerationSet(Description, Request, HTTPResponseCode, Path, Query, Paste, Tags),
	HTTPResponseCode: createEnumerationSet(Body, Headers, Paste),
	Request:          createEnumerationSet(Body, Headers, Paste),
	Info:             createEnumerationSet(Title, Version, Description, Paste),
	Server:           createEnumerationSet(BaseUrl, Paste),
	Method:           createEnumerationSet(Description, Params, Result, Tags),
	TAG:              createEnumerationSet(Description),
	Macro: createEnumerationSet(
		Info,
		Title,
		Version,
		Description,
		Server,
		BaseUrl,
		Url,
		Get,
		Post,
		Put,
		Patch,
		Delete,
		Body,
		Request,
		HTTPResponseCode,
		Path,
		Headers,
		Query,
		Type,
		Enum,
		Paste,
	),
}

func createEnumerationSet(ee ...Enumeration) map[Enumeration]struct{} {
	if len(ee) == 0 {
		return nil
	}

	res := make(map[Enumeration]struct{}, len(ee))
	for _, e := range ee {
		res[e] = struct{}{}
	}
	return res
}

func IsStartWithDirective(b bytes.Bytes) bool {
	if len(b) < 3 {
		return false
	}

	switch b[0] { // response directive 100, 200, 300 etc
	case '1', '2', '3', '4', '5':
		if IsHTTPResponseCode(string(b[0:3])) {
			return true
		}
	}

	s := string(b)

	for i := 0; i < len(ss); i++ {
		de := Enumeration(i)
		if de == HTTPResponseCode {
			continue
		}
		if strings.HasPrefix(s, de.String()) {
			return true
		}
	}

	return false
}
