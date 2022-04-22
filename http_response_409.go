package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func httpResponse409(w http.ResponseWriter, e error) {
	info := newErrorInfo(e)
	b, err := json.Marshal(info)
	if err != nil {
		httpResponse500(w, err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	_, _ = w.Write(b)

	log.Print("... " + e.Error())
}
