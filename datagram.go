package main

import (
	"log"
	"strconv"

	"github.com/jsightapi/datagram"
	"github.com/jsightapi/jsight-api-core/jerr"
	"github.com/jsightapi/jsight-api-core/kit"
)

func sendDatagram(clientID, clientIP string, projectSize int, j kit.JApi, je *jerr.JApiError) {
	d := datagram.New()
	d.Append("cid", clientID)
	d.Append("cip", clientIP)
	d.Append("at", "1")                       // Application Type
	d.AppendTruncatable("pt", j.Title())      // Project title
	d.Append("ps", strconv.Itoa(projectSize)) // Project size
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
