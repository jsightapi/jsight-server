package catalog

import (
	"sync"
)

// Tags represent available tags.
// gen:OrderedMap
type Tags struct {
	data  map[TagName]*Tag
	order []TagName
	mx    sync.RWMutex
}
