package directive

import (
	"sync"
)

// Directives represent available directives.
// gen:OrderedMap
type Directives struct {
	data  map[string]*Directive
	order []string
	mx    sync.RWMutex
}
