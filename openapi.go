package main

import (
	"bytes"
	"encoding/json"

	"github.com/itchyny/json2yaml"

	"github.com/jsightapi/jsight-api-core/catalog/ser/openapi"
	"github.com/jsightapi/jsight-api-core/kit"
)

func openapiJSON(jAPI kit.JApi) ([]byte, error) {
	oa, oaErr := openapi.NewOpenAPI(jAPI.Catalog())
	if oaErr != nil {
		return nil, oaErr
	}

	resp, jsonErr := json.MarshalIndent(oa, "", "  ")
	if jsonErr != nil {
		return nil, jsonErr
	}

	return resp, nil
}

func openapiYAML(jAPI kit.JApi) ([]byte, error) {
	js, err := openapiJSON(jAPI)
	if err != nil {
		return nil, err
	}

	resp, err := jsonToYAML(js)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func jsonToYAML(jsonData []byte) ([]byte, error) {
	from := bytes.NewReader(jsonData)
	to := bytes.NewBuffer(make([]byte, 0, len(jsonData)))

	err := json2yaml.Convert(to, from)
	if err != nil {
		return nil, err
	}

	return removeLastLineBreak(to), nil
}

func removeLastLineBreak(input *bytes.Buffer) []byte {
	if bs := input.Bytes(); len(bs) > 0 && bs[len(bs)-1] == '\n' {
		return bs[:len(bs)-1]
	}
	return input.Bytes()
}
