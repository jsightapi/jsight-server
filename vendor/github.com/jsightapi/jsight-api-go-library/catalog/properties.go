package catalog

import (
	"sync"
)

// Properties represent JSight object properties.
// gen:OrderedMap
type Properties struct {
	data  map[string]*SchemaContentJSight
	order []string
	mx    sync.RWMutex
}

func NewProperties(data map[string]*SchemaContentJSight, order []string) *Properties {
	return &Properties{
		data:  data,
		order: order,
	}
}
