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
