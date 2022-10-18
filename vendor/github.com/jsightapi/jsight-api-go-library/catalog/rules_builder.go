package catalog

import "sync"

type RulesBuilder struct {
	rules *Rules
	mx    sync.RWMutex
}

func newRulesBuilder(caption int) *RulesBuilder {
	return &RulesBuilder{
		rules: &Rules{
			data:  make([]Rule, 0, caption),
			index: make(map[string]int, caption),
		},
	}
}

func (b *RulesBuilder) Set(k string, r Rule) {
	b.mx.Lock()
	defer b.mx.Unlock()

	r.Key = k

	b.rules.index[k] = len(b.rules.data)
	b.rules.data = append(b.rules.data, r)
}

func (b *RulesBuilder) Append(r Rule) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.rules.data = append(b.rules.data, r)
}

func (b *RulesBuilder) Rules() *Rules {
	return b.rules
}
