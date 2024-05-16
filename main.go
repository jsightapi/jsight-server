package main

import (
	_ "embed"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/convert-jsight", convertJSight)
	// http.HandleFunc("/", pageReload)

	server := &http.Server{
		Addr:        ":8080",
		ReadTimeout: 5 * time.Second,
	}

	log.Print("The server is running on the port :8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
