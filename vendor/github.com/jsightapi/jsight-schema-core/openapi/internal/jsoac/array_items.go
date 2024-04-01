package jsoac

import (
	"encoding/json"
)

type ArrayItems struct {
	items []ArrayItem
}

type ArrayItem struct {
	value Node
}

var _ json.Marshaler = ArrayItems{}
var _ json.Marshaler = &ArrayItems{}

func newArrayItems(length int) ArrayItems {
	return ArrayItems{items: make([]ArrayItem, 0, length)}
}

func (ai *ArrayItems) append(value Node) {
	i := ArrayItem{
		value: value,
	}
	ai.items = append(ai.items, i)
}

func (ai ArrayItems) MarshalJSON() ([]byte, error) {
	b := bufferPool.Get()
	defer bufferPool.Put(b)
	length := len(ai.items)

	if length > 1 {
		b.WriteString(`{"anyOf": [`)
	}
	for i, item := range ai.items {
		value, err := json.Marshal(item.value)
		if err != nil {
			return nil, err
		}
		b.Write(value)

		if i+1 != length {
			b.WriteByte(',')
		}
	}
	if length > 1 {
		b.WriteString(`]}`)
	}
	if length == 0 {
		b.WriteString(`{}`)
	}
	return b.Bytes(), nil
}
