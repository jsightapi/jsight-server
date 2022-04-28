package catalog

import (
	"sync"
)

// UserTypes represent available user types.
// gen:OrderedMap
type UserTypes struct {
	data  map[string]*UserType
	order []string
	mx    sync.RWMutex
}
