package main

import (
	"errors"
	"j/japi/jerr"
)

type errorInfo struct {
	Status  string
	Message string
	Line    int
	Index   int
}

func newErrorInfo(e error) errorInfo {
	r := errorInfo{
		Status:  "Error",
		Message: e.Error(),
	}

	var je *jerr.JAPIError

	if errors.As(e, &je) {
		r.Line = int(je.Line())
		r.Index = int(je.Index())
	}

	return r
}
