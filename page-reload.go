package main

import (
	"log"
	"net/http"
)

func pageReload(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	format := r.FormValue("format")
	log.Printf("%s %s %s %s", r.Method, r.URL.Path, to, format)

	if getBoolEnv("JSIGHT_SERVER_CORS") {
		cors(w)
	}

	wr := httpResponseWriter{writer: w}

	switch r.Method {
	case http.MethodOptions:

	case http.MethodPost:
		pageReloadPOST(wr, r)
		return

	default:
		wr.errorStr("HTTP POST request required")
		return
	}
}

func pageReloadPOST(wr httpResponseWriter, r *http.Request) {
	wr.errorPageReload()
	return
}
