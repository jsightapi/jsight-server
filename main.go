package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jsightapi/datagram"
	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-api-core/kit"
	"github.com/jsightapi/jsight-schema-core/fs"
)

func main() {
	http.HandleFunc("/", jdocExchangeFile)

	server := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
	}

	log.Print("The server is running on the port :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func jdocExchangeFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	if getBoolEnv("JSIGHT_SERVER_CORS") {
		cors(w)
	}

	switch r.Method {
	case http.MethodOptions:

	case http.MethodPost:
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

		httpResponse200(w, json)

	default:
		httpResponse409(w, errors.New("HTTP POST request required"))
	}
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
