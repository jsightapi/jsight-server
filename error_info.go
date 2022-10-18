package main

import (
	"errors"

	"github.com/jsightapi/jsight-api-go-library/jerr"
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

	var je *jerr.JApiError

	if errors.As(e, &je) {
		r.Line = int(je.Line())
		r.Index = int(je.Index())
	}

	return r
}
