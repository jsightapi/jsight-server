package catalog

import (
	"sync"
)

// HttpInteractions represent available resource methods.
// gen:OrderedMap
type HttpInteractions struct {
	data  map[HttpInteractionId]*HttpInteraction
	order []HttpInteractionId
	mx    sync.RWMutex
}
