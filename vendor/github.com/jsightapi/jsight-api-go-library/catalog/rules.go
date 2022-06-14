package catalog

import (
	"sync"
)

// Rules represent JSight rules.
// gen:OrderedMap
type Rules struct {
	data  map[string]Rule
	order []string
	mx    sync.RWMutex
}

func NewRules(data map[string]Rule, order []string) *Rules {
	return &Rules{
		data:  data,
		order: order,
	}
}
