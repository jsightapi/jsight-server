package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jsightapi/datagram"
	"github.com/jsightapi/jsight-api-go-library/kit"
	"github.com/jsightapi/jsight-schema-go-library/fs"
)

func main() {
	http.HandleFunc("/", jdocExchangeFile)

	log.Print("The server is running on the port :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
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

		j := kit.NewJApiFromFile(fs.NewFile("root", b))
		je := j.ValidateJAPI()

		if getBoolEnv("JSIGHT_SERVER_STATISTICS") {
			d := datagram.New()
			d.Append("cid", r.Header.Get("X-Browser-UUID")) // Client ID
			d.Append("cip", getIP(r))                       // Client IP
			d.Append("at", "1")                             // Application Type
			d.AppendTruncatable("pt", j.Title())            // Project title
			d.Append("ps", strconv.Itoa(len(b)))            // Project size
			if je != nil {
				d.AppendTruncatable("pem", je.Error())                      // Project error message
				d.Append("pel", strconv.FormatUint(uint64(je.Line()), 10))  // Project error line
				d.Append("pei", strconv.FormatUint(uint64(je.Index()), 10)) // Project error index
			}

			err = sendToStatisticServer(d.Pack())
			if err != nil {
				log.Print("... " + err.Error())
			}
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
