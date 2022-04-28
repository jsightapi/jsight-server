package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/jsightapi/jsight-api-go-library/kit"
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

		j := kit.NewJapiFromBytes(b)
		je := j.ValidateJAPI()

		if getBoolEnv("JSIGHT_SERVER_STATISTICS") {
			e := event{
				ClientIPv4:      getIP(r),
				SourceType:      "Editor",
				SourceID:        r.Header.Get("X-Browser-UUID"),
				ProjectTitle:    j.Title(),
				ProjectSize:     uint32(len(b)),
				ProjectHasError: je != nil,
			}
			err = sendToStatisticServer(e)
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
