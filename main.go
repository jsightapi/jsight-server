package main

import (
	_ "embed"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jsightapi/jsight-schema-core/fs"

	"github.com/jsightapi/datagram"
	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-api-core/kit"
)

//go:embed testdata/openapi.json
var openapiJSON []byte

//go:embed testdata/openapi.yaml
var openapiYAML []byte

func main() {
	http.HandleFunc("/", convertJSight)

	server := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
	}

	log.Print("The server is running on the port :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func convertJSight(w http.ResponseWriter, r *http.Request) {
	to := r.FormValue("to")
	format := r.FormValue("format")

	log.Printf("%s %s %s %s", r.Method, r.URL.Path, to, format)

	if getBoolEnv("JSIGHT_SERVER_CORS") {
		cors(w)
	}

	switch r.Method {
	case http.MethodOptions:

	case http.MethodPost:
		switch to {
		case "jdoc-2.0":
			convertToJDoc(w, r)
		case "openapi-3.0.3":
			convertToOpenAPI(w, r)
		default:
			httpResponse409(w, errors.New("not supported"))
		}
	default:
		httpResponse409(w, errors.New("HTTP POST request required"))
	}
}

func convertToOpenAPI(w http.ResponseWriter, r *http.Request) {
	switch r.FormValue("format") {
	case "json", "":
		convertToOpenapiJSON(w)
	case "yaml":
		convertToOpenapiYAML(w)
	default:
		httpResponse409(w, errors.New("not supported format"))
	}
}

func convertToOpenapiJSON(w http.ResponseWriter) {
	httpResponseJSON200(w, openapiJSON)
}

func convertToOpenapiYAML(w http.ResponseWriter) {
	httpResponseYAML200(w, openapiYAML)
}

func convertToJDoc(w http.ResponseWriter, r *http.Request) {
	format := r.FormValue("format")
	if format != "json" && format != "" {
		httpResponse409(w, errors.New("not supported format"))
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		httpResponse409(w, err)
		return
	}

	j, je := kit.NewJApiFromFile(fs.NewFile("root", b))

	if getBoolEnv("JSIGHT_SERVER_STATISTICS") {
		sendDatagram(r, len(b), j, je)
	}

	if je != nil {
		httpResponse409(w, je)
		return
	}

	json, err := j.ToJson()
	if err != nil {
		httpResponse409(w, err)
		return
	}

	httpResponseJDoc200(w, json)
}

func sendDatagram(r *http.Request, projectSize int, j kit.JApi, je *jerr.JApiError) {
	d := datagram.New()
	d.Append("cid", r.Header.Get("X-Browser-UUID")) // Client ID
	d.Append("cip", getIP(r))                       // Client IP
	d.Append("at", "1")                             // Application Type
	d.AppendTruncatable("pt", j.Title())            // Project title
	d.Append("ps", strconv.Itoa(projectSize))       // Project size
	if je != nil {
		d.AppendTruncatable("pem", je.Error())                    // Project error message
		d.Append("pel", strconv.FormatUint(uint64(je.Line), 10))  // Project error line
		d.Append("pei", strconv.FormatUint(uint64(je.Index), 10)) // Project error index
	}

	err := sendToStatisticServer(d.Pack())
	if err != nil {
		log.Print("... " + err.Error())
	}
}
