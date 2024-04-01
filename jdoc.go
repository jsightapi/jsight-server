package main

import (
	"github.com/jsightapi/jsight-api-core/kit"
)

func jdocJSON(jAPI kit.JApi) ([]byte, error) {
	json, err := jAPI.ToJson()
	if err != nil {
		return nil, err
	}

	return json, nil
}
