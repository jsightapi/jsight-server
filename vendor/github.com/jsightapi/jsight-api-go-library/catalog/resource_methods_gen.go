// Autogenerated code!
// DO NOT EDIT!
//
// Generated by OrderedMap generator from the internal/cmd/generator command.

package catalog

import (
	"bytes"
	"encoding/json"
)

// Set sets a value with specified key.
func (m *ResourceMethods) Set(k ResourceMethodId, v *ResourceMethod) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[ResourceMethodId]*ResourceMethod{}
	}
	if !m.has(k) {
		m.order = append(m.order, k)
	}
	m.data[k] = v
}

// SetToTop do the same as Set, but new key will be placed on top of the order
// map.
func (m *ResourceMethods) SetToTop(k ResourceMethodId, v *ResourceMethod) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[ResourceMethodId]*ResourceMethod{}
	}
	if !m.has(k) {
		m.order = append([]ResourceMethodId{k}, m.order...)
	}
	m.data[k] = v
}

// Update updates a value with specified key.
func (m *ResourceMethods) Update(k ResourceMethodId, fn func(v *ResourceMethod) *ResourceMethod) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if !m.has(k) {
		// Prevent from possible nil pointer dereference if map value type is a
		// pointer.
		return
	}

	m.data[k] = fn(m.data[k])
}

// GetValue gets a value by key.
func (m *ResourceMethods) GetValue(k ResourceMethodId) *ResourceMethod {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.data[k]
}

// Get gets a value by key.
func (m *ResourceMethods) Get(k ResourceMethodId) (*ResourceMethod, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	v, ok := m.data[k]
	return v, ok
}

// Has checks that specified key is set.
func (m *ResourceMethods) Has(k ResourceMethodId) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.has(k)
}

func (m *ResourceMethods) has(k ResourceMethodId) bool {
	_, ok := m.data[k]
	return ok
}

// Len returns count of values.
func (m *ResourceMethods) Len() int {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return len(m.data)
}

// Find finds first matched item from the map.
func (m *ResourceMethods) Find(fn findResourceMethodsFunc) (ResourceMethodsItem, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		if fn(k, m.data[k]) {
			return ResourceMethodsItem{
				Key:   k,
				Value: m.data[k],
			}, true
		}
	}
	return ResourceMethodsItem{}, false
}

type findResourceMethodsFunc = func(k ResourceMethodId, v *ResourceMethod) bool

// Each iterates and perform given function on each item in the map.
func (m *ResourceMethods) Each(fn eachResourceMethodsFunc) error {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		if err := fn(k, m.data[k]); err != nil {
			return err
		}
	}
	return nil
}

// EachReverse act almost the same as Each but in reverse order.
func (m *ResourceMethods) EachReverse(fn eachResourceMethodsFunc) error {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for i := len(m.order) - 1; i >= 0; i-- {
		k := m.order[i]
		if err := fn(k, m.data[k]); err != nil {
			return err
		}
	}
	return nil
}

type eachResourceMethodsFunc = func(k ResourceMethodId, v *ResourceMethod) error

func (m *ResourceMethods) EachSafe(fn eachSafeResourceMethodsFunc) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		fn(k, m.data[k])
	}
}

type eachSafeResourceMethodsFunc = func(k ResourceMethodId, v *ResourceMethod)

// Map iterates and changes values in the map.
func (m *ResourceMethods) Map(fn mapResourceMethodsFunc) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	for _, k := range m.order {
		v, err := fn(k, m.data[k])
		if err != nil {
			return err
		}
		m.data[k] = v
	}
	return nil
}

type mapResourceMethodsFunc = func(k ResourceMethodId, v *ResourceMethod) (*ResourceMethod, error)

// ResourceMethodsItem represent single data from the ResourceMethods.
type ResourceMethodsItem struct {
	Key   ResourceMethodId
	Value *ResourceMethod
}

var _ json.Marshaler = &ResourceMethods{}

func (m *ResourceMethods) MarshalJSON() ([]byte, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	var buf bytes.Buffer
	buf.WriteRune('{')

	for i, k := range m.order {
		if i != 0 {
			buf.WriteRune(',')
		}

		// marshal key
		key, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteRune(':')

		// marshal value
		val, err := json.Marshal(m.data[k])
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}
