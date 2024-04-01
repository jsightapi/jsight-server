package main

import (
	"encoding/json"

	"github.com/jsightapi/jsight-api-core/catalog/ser/openapi"
	"github.com/jsightapi/jsight-api-core/kit"
	"gopkg.in/yaml.v3"
)

func openapiJSON(jAPI kit.JApi) ([]byte, error) {
	oa, err := openapi.NewOpenAPI(jAPI.Catalog())
	if err != nil {
		return nil, err
	}

	resp, err := json.Marshal(oa)
	if err != nil {
		return nil, err
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

func jsonToYAML(j []byte) ([]byte, error) {
	var data any

	err := yaml.Unmarshal(j, &data)
	if err != nil {
		return nil, err
	}

	return yaml.Marshal(data)
}
