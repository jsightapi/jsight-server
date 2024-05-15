package catalog

import (
	"sync"
)

// UserRules represent available user rules.
// gen:OrderedMap
type UserRules struct {
	data  map[string]*UserRule
	order []string
	mx    sync.RWMutex
}
