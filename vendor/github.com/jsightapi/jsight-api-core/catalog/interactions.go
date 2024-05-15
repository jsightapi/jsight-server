package catalog

import (
	"sync"
)

// Interactions represent available resource methods.
// gen:OrderedMap
type Interactions struct {
	data  map[InteractionID]Interaction
	order []InteractionID
	mx    sync.RWMutex
}
