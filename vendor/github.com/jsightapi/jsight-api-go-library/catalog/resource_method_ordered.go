package catalog

import (
	"bytes"
	"encoding/json"
)

type Resource struct {
	Val *ResourceMethod
	Key ResourceMethodId
}

type OrderedResources []Resource

func (oo OrderedResources) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteRune('{')

	for i, kv := range oo {
		if i != 0 {
			buf.WriteRune(',')
		}

		// marshal key
		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteRune(':')

		// marshal value
		val, err := json.Marshal(kv.Val)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}
