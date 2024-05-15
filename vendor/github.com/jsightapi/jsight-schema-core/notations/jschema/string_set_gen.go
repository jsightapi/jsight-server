// Autogenerated code!
// DO NOT EDIT!
//
// Generated by Set generator from the internal/cmd/generator command.

package jschema

func NewStringSet(vv ...string) *StringSet {
	data := make(map[string]struct{}, len(vv))

	for _, v := range vv {
		data[v] = struct{}{}
	}

	return &StringSet{
		data:  data,
		order: vv,
	}
}

// Add adds specific value to set.
func (m *StringSet) Add(v string) {
	m.mx.Lock()
	defer m.mx.Unlock()

	if m.data == nil {
		m.data = map[string]struct{}{}
	}
	if !m.has(v) {
		m.order = append(m.order, v)
	}
	m.data[v] = struct{}{}
}

// Has checks that specified value is exists.
func (m *StringSet) Has(v string) bool {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.has(v)
}

func (m *StringSet) has(v string) bool {
	_, ok := m.data[v]
	return ok
}

// Len returns len of set.
func (m *StringSet) Len() int {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return len(m.data)
}

// Data return set's data.
func (m *StringSet) Data() []string {
	m.mx.RLock()
	defer m.mx.RUnlock()

	return m.order
}
