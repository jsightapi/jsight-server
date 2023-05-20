package main

import (
	"errors"

	"github.com/jsightapi/jsight-api-core/jerr"
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
		r.Line = je.Line.Int()
		r.Index = je.Index.Int()
	}

	return r
}
