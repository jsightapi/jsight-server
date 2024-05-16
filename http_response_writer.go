package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jsightapi/jsight-api-core/catalog"
)

type httpResponseWriter struct {
	writer http.ResponseWriter
}

func (r httpResponseWriter) jdocJSON(b []byte) {
	r.writer.Header().Set("X-Jdoc-Exchange-Version", catalog.JDocExchangeVersion)
	r.json(b)
}

func (r httpResponseWriter) json(b []byte) {
	r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	n, _ := r.writer.Write(b)

	log.Printf("... Ok (%d bytes)", n)
}

func (r httpResponseWriter) yaml(b []byte) {
	r.writer.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	n, _ := r.writer.Write(b)

	log.Printf("... Ok (%d bytes)", n)
}

func (r httpResponseWriter) errorStr(s string) {
	r.error(errors.New(s))
}

func (r httpResponseWriter) error(e error) {
	info := newErrorInfo(e)
	b, err := json.Marshal(info)
	if err != nil {
		r.internalServerError(err)
		return
	}

	r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.writer.WriteHeader(http.StatusConflict)
	_, _ = r.writer.Write(b)

	log.Print("... " + e.Error())
}

func (r httpResponseWriter) internalServerError(e error) {
	r.writer.Header().Set("Content-Type", "text/plain")
	r.writer.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprint(r.writer, e.Error())

	log.Print("... " + e.Error())
}

func (r httpResponseWriter) errorPageReload() {
	info := errorInfo{
		"Error",
		"Please hard refresh your browser",
		0, 0,
	}

	b, err := json.Marshal(info)
	if err != nil {
		r.internalServerError(err)
		return
	}

	r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	r.writer.WriteHeader(http.StatusConflict)
	_, _ = r.writer.Write(b)

	log.Print("... " + "returned error message for page reloading")
}
