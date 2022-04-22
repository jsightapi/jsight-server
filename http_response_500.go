package main

import (
	"fmt"
	"log"
	"net/http"
)

func httpResponse500(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprint(w, e.Error())

	log.Print("... " + e.Error())
}
