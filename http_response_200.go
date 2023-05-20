package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jsightapi/jsight-api-core/catalog"
)

func httpResponse200(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Jdoc-Exchange-Version", catalog.JDocExchangeVersion)
	n, _ := fmt.Fprint(w, string(b))

	log.Printf("... Ok (%d bytes)", n)
}
