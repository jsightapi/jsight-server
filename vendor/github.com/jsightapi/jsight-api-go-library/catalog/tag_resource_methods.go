package catalog

import (
	"sync"
)

// TagResourceMethods represent available resource methods for the tag.
// gen:OrderedMap
type TagResourceMethods struct {
	data  map[Path]*ResourceMethodIdList
	order []Path
	mx    sync.RWMutex
}
