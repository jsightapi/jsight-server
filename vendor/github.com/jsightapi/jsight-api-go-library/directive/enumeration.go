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
)

var (
	ss = [...]string{
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
	switch de {
	case Get, Post, Put, Patch, Delete:
		return true
	default:
		return false
	}
}

func IsAllowedForRootContext(child Enumeration) bool {
	switch child {
	case Jsight, Url, Info, Server, Type, Enum,
		Get, Post, Put, Patch, Delete, Macro, Paste:
		return true
	default:
		return false
	}
}

func (de Enumeration) IsAllowedForDirectiveContext(child Enumeration) bool {
	switch de {
	case Url:
		switch child {
		case Path, Paste,
			Get, Post, Put, Patch, Delete:
			return true
		default:
			return false
		}

	case Get, Post, Put, Patch, Delete:
		switch child {
		case Description, Query, Path, Request, HTTPResponseCode, Paste:
			return true
		default:
			return false
		}

	case HTTPResponseCode, Request:
		switch child {
		case Body, Headers, Paste:
			return true
		default:
			return false
		}

	case Info:
		switch child {
		case Title, Version, Description, Paste:
			return true
		default:
			return false
		}

	case Server:
		switch child {
		case BaseUrl, Paste:
			return true
		default:
			return false
		}

	case Macro:
		switch child {
		case Info, Title, Version, Description, Server, BaseUrl, Url, Get, Post, Put, Patch, Delete, Body,
			Request, HTTPResponseCode, Path, Headers, Query, Type, Enum, Paste:
			return true
		default:
			return false
		}

	default: // Body, Query, Path, Description, Jsight, Type, Title, Version, BaseUrl, Paste
		return false
	}
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
