// Autogenerated code!
// DO NOT EDIT!
//
// Generated by UnsafeOrderedMap generator from the internal/cmd/generator command.

package catalog

import (
	"bytes"
	"encoding/json"

	schema "github.com/jsightapi/jsight-schema-core"
)

// Set sets a value with specified key.
func (m *UserSchemas) Set(k string, v schema.Schema) {
	if m.data == nil {
		m.data = map[string]schema.Schema{}
	}
	if !m.Has(k) {
		m.order = append(m.order, k)
	}
	m.data[k] = v
}

// Update updates a value with specified key.
func (m *UserSchemas) Update(k string, fn func(v schema.Schema) schema.Schema) {
	if !m.Has(k) {
		// Prevent from possible nil pointer dereference if map value type is a
		// pointer.
		return
	}

	m.data[k] = fn(m.data[k])
}

// GetValue gets a value by key.
func (m *UserSchemas) GetValue(k string) schema.Schema {
	return m.data[k]
}

// Get gets a value by key.
func (m *UserSchemas) Get(k string) (schema.Schema, bool) {
	v, ok := m.data[k]
	return v, ok
}

// Has checks that specified key is set.
func (m *UserSchemas) Has(k string) bool {
	_, ok := m.data[k]
	return ok
}

// Len returns count of values.
func (m *UserSchemas) Len() int {
	return len(m.data)
}

// Each iterates and perform given function on each item in the map.
func (m *UserSchemas) Each(fn eachUserSchemasFunc) error {
	for _, k := range m.order {
		if err := fn(k, m.data[k]); err != nil {
			return err
		}
	}
	return nil
}

type eachUserSchemasFunc = func(k string, v schema.Schema) error

func (m *UserSchemas) EachSafe(fn eachSafeUserSchemasFunc) {
	for _, k := range m.order {
		fn(k, m.data[k])
	}
}

type eachSafeUserSchemasFunc = func(k string, v schema.Schema)

// Map iterates and changes values in the map.
func (m *UserSchemas) Map(fn mapUserSchemasFunc) error {
	for _, k := range m.order {
		v, err := fn(k, m.data[k])
		if err != nil {
			return err
		}
		m.data[k] = v
	}
	return nil
}

type mapUserSchemasFunc = func(k string, v schema.Schema) (schema.Schema, error)

// UserSchemasItem represent single data from the UserSchemas.
type UserSchemasItem struct {
	Key   string
	Value schema.Schema
}

var _ json.Marshaler = &UserSchemas{}

func (m *UserSchemas) MarshalJSON() ([]byte, error) {
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