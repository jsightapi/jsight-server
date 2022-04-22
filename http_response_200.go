package main

import (
	"fmt"
	"j/japi/catalog"
	"log"
	"net/http"
)

func httpResponse200(w http.ResponseWriter, b []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Jdoc-Exchange-File-Schema-Version", catalog.JDocExchangeFileSchemaVersion)
	n, _ := fmt.Fprint(w, string(b))

	log.Printf("... Ok (%d bytes)", n)
}
