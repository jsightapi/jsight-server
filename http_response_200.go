package main

import (
	"log"
	"net/http"

	"github.com/jsightapi/jsight-api-core/catalog"
)

func httpResponseJDoc200(w http.ResponseWriter, b []byte) {
	w.Header().Set("X-Jdoc-Exchange-Version", catalog.JDocExchangeVersion)
	httpResponseJSON200(w, b)
}

func httpResponseJSON200(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	n, _ := w.Write(b)

	log.Printf("... Ok (%d bytes)", n)
}

func httpResponseYAML200(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	n, _ := w.Write(b)

	log.Printf("... Ok (%d bytes)", n)
}
