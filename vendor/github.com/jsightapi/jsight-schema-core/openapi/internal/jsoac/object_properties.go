package jsoac

import (
	"encoding/json"
)

type ObjectProperties struct {
	properties []Property
}

type Property struct {
	key   string
	value Node
}

var _ json.Marshaler = ObjectProperties{}
var _ json.Marshaler = &ObjectProperties{}

func newObjectProperties(length int) ObjectProperties {
	return ObjectProperties{properties: make([]Property, 0, length)}
}

func (op *ObjectProperties) append(key string, value Node) {
	p := Property{
		key:   key,
		value: value,
	}
	op.properties = append(op.properties, p)
}

func (op ObjectProperties) MarshalJSON() ([]byte, error) {
	b := bufferPool.Get()
	defer bufferPool.Put(b)

	b.WriteByte('{')
	length := len(op.properties)
	for i, property := range op.properties {
		value, err := json.Marshal(property.key)
		if err != nil {
			return nil, err
		}
		b.Write(value)
		b.WriteString(`:`)

		value, err = json.Marshal(property.value)
		if err != nil {
			return nil, err
		}
		b.Write(value)

		if i+1 != length {
			b.WriteByte(',')
		}
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}
