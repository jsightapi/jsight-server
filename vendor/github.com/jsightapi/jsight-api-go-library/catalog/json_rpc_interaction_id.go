package catalog

import (
	"fmt"

	"github.com/jsightapi/jsight-api-go-library/directive"
)

type JsonRpcInteractionId struct {
	protocol Protocol
	path     Path
	method   string
}

func (j JsonRpcInteractionId) Protocol() Protocol {
	return j.protocol
}

func (j JsonRpcInteractionId) Path() Path {
	return j.path
}

func (j JsonRpcInteractionId) String() string {
	return fmt.Sprintf("json-rpc-2.0 %s %s", j.method, j.path.String())
}

func (j JsonRpcInteractionId) MarshalText() ([]byte, error) {
	return []byte(j.String()), nil
}

func newJsonRpcInteractionId(d directive.Directive) (JsonRpcInteractionId, error) {
	j := JsonRpcInteractionId{
		protocol: JsonRpc,
	}

	path, err := d.Path()
	if err != nil {
		return j, err
	}

	j.path = Path(path)

	j.method, err = d.JsonRpcMethodName()
	if err != nil {
		return j, err
	}

	return j, nil
}
