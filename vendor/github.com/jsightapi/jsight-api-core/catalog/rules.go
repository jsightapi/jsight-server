package catalog

import (
	"encoding/json"
)

type Rules struct {
	index map[string]int
	data  []Rule
}

var _ json.Marshaler = &Rules{}

func NewRules(d []Rule) *Rules {
	rr := &Rules{
		data:  d,
		index: make(map[string]int, len(d)),
	}
	for i, r := range d {
		rr.index[r.Key] = i
	}
	return rr
}

func (rr *Rules) Len() int {
	return len(rr.data)
}

func (rr *Rules) Has(k string) bool {
	if rr == nil {
		return false
	}
	_, ok := rr.index[k]
	return ok
}

func (rr *Rules) Get(k string) (Rule, bool) {
	if rr == nil {
		return Rule{}, false
	}

	i, ok := rr.index[k]
	if !ok {
		return Rule{}, false
	}

	return rr.data[i], true
}

type findRulesFunc = func(k string, v Rule) error

func (rr *Rules) Each(fn findRulesFunc) error {
	if rr != nil {
		for _, v := range rr.data {
			if err := fn(v.Key, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (rr *Rules) MarshalJSON() ([]byte, error) {
	return json.Marshal(rr.data)
}
