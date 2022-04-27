package catalog

import (
	"sync"
)

// ResourceMethods represent available resource methods.
// gen:OrderedMap
type ResourceMethods struct {
	data  map[ResourceMethodId]*ResourceMethod
	order []ResourceMethodId
	mx    sync.RWMutex
}
