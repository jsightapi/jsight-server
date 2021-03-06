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
func (m *TagResourceMethods) Set(k Path, v *ResourceMethodIdList) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[Path]*ResourceMethodIdList{}
	}
	if !m.has(k) {
		m.order = append(m.order, k)
	}
	m.data[k] = v
}

// SetToTop do the same as Set, but new key will be placed on top of the order
// map.
func (m *TagResourceMethods) SetToTop(k Path, v *ResourceMethodIdList) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[Path]*ResourceMethodIdList{}
	}
	if !m.has(k) {
		m.order = append([]Path{k}, m.order...)
	}
	m.data[k] = v
}

// Update updates a value with specified key.
func (m *TagResourceMethods) Update(k Path, fn func(v *ResourceMethodIdList) *ResourceMethodIdList) {
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
func (m *TagResourceMethods) GetValue(k Path) *ResourceMethodIdList {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.data[k]
}

// Get gets a value by key.
func (m *TagResourceMethods) Get(k Path) (*ResourceMethodIdList, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	v, ok := m.data[k]
	return v, ok
}

// Has checks that specified key is set.
func (m *TagResourceMethods) Has(k Path) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.has(k)
}

func (m *TagResourceMethods) has(k Path) bool {
	_, ok := m.data[k]
	return ok
}

// Len returns count of values.
func (m *TagResourceMethods) Len() int {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return len(m.data)
}

// Find finds first matched item from the map.
func (m *TagResourceMethods) Find(fn findTagResourceMethodsFunc) (TagResourceMethodsItem, bool) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		if fn(k, m.data[k]) {
			return TagResourceMethodsItem{
				Key:   k,
				Value: m.data[k],
			}, true
		}
	}
	return TagResourceMethodsItem{}, false
}

type findTagResourceMethodsFunc = func(k Path, v *ResourceMethodIdList) bool

// Each iterates and perform given function on each item in the map.
func (m *TagResourceMethods) Each(fn eachTagResourceMethodsFunc) error {
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
func (m *TagResourceMethods) EachReverse(fn eachTagResourceMethodsFunc) error {
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

type eachTagResourceMethodsFunc = func(k Path, v *ResourceMethodIdList) error

func (m *TagResourceMethods) EachSafe(fn eachSafeTagResourceMethodsFunc) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	for _, k := range m.order {
		fn(k, m.data[k])
	}
}

type eachSafeTagResourceMethodsFunc = func(k Path, v *ResourceMethodIdList)

// Map iterates and changes values in the map.
func (m *TagResourceMethods) Map(fn mapTagResourceMethodsFunc) error {
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

type mapTagResourceMethodsFunc = func(k Path, v *ResourceMethodIdList) (*ResourceMethodIdList, error)

// TagResourceMethodsItem represent single data from the TagResourceMethods.
type TagResourceMethodsItem struct {
	Key   Path
	Value *ResourceMethodIdList
}

var _ json.Marshaler = &TagResourceMethods{}

func (m *TagResourceMethods) MarshalJSON() ([]byte, error) {
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
