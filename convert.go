package main

import (
	"io"
	"log"
	"net/http"

	"github.com/jsightapi/jsight-schema-core/fs"

	"github.com/jsightapi/jsight-api-core/kit"
)

func convertJSight(w http.ResponseWriter, r *http.Request) {
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
		convertJSightPOST(wr, r)
		return

	default:
		wr.errorStr("HTTP POST request required")
		return
	}
}

func convertJSightPOST(wr httpResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	format := r.FormValue("format")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		wr.error(err)
		return
	}

	jAPI, jErr := kit.NewJApiFromFile(fs.NewFile("root", body))

	if getBoolEnv("JSIGHT_SERVER_STATISTICS") {
		clientID := r.Header.Get("X-Browser-UUID")
		clientIP := getIP(r)
		sendDatagram(clientID, clientIP, len(body), jAPI, jErr)
	}

	if jErr != nil {
		wr.error(jErr)
		return
	}

	switch to {
	case "jdoc-2.0":
		switch format {
		case "json", "":
			writeJDocJSON(wr, jAPI)
			return
		default:
			wr.errorStr("not supported format")
			return
		}
	case "openapi-3.0.3":
		switch format {
		case "json", "":
			writeOpenapiJSON(wr, jAPI)
			return
		case "yaml":
			writeOpenapiYAML(wr, jAPI)
			return
		default:
			wr.errorStr("not supported format")
			return
		}
	default:
		wr.errorStr(`you must specify the "to" parameter`)
		return
	}
}

func writeJDocJSON(wr httpResponseWriter, jAPI kit.JApi) {
	resp, err := jdocJSON(jAPI)
	if err != nil {
		wr.error(err)
		return
	}

	wr.jdocJSON(resp)
}

func writeOpenapiJSON(wr httpResponseWriter, jAPI kit.JApi) {
	resp, err := openapiJSON(jAPI)
	if err != nil {
		wr.error(err)
		return
	}

	wr.json(resp)
}

func writeOpenapiYAML(wr httpResponseWriter, jAPI kit.JApi) {
	resp, err := openapiYAML(jAPI)
	if err != nil {
		wr.error(err)
		return
	}

	wr.yaml(resp)
}
