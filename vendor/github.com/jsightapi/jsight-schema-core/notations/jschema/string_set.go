package jschema

import "sync"

// StringSet a set of strings.
// gen:Set
type StringSet struct {
	data  map[string]struct{}
	order []string
	mx    sync.RWMutex
}
