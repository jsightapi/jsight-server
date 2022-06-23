package schema

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

type ObjectNodeKeys struct {
	index map[string]int
	Data  []ObjectNodeKey
}

type ObjectNodeKey struct {
	Key        string
	Lex        lexeme.LexEvent
	Index      int
	IsShortcut bool
}

func newObjectNodeKeys() *ObjectNodeKeys {
	return &ObjectNodeKeys{
		Data:  make([]ObjectNodeKey, 0, 5),
		index: make(map[string]int, 5),
	}
}

func (k *ObjectNodeKeys) Set(v ObjectNodeKey) {
	if !v.IsShortcut {
		if _, ok := k.index[v.Key]; ok {
			panic(errors.Format(errors.ErrDuplicateKeysInSchema, v.Key))
		}
	}

	k.index[v.Key] = v.Index
	k.Data = append(k.Data, v)
}

func (k ObjectNodeKeys) Find(i int) (ObjectNodeKey, bool) {
	if len(k.Data) > i {
		return k.Data[i], true
	}
	return ObjectNodeKey{}, false
}

func (k ObjectNodeKeys) Get(key string) (ObjectNodeKey, bool) {
	if i, ok := k.index[key]; ok {
		return k.Data[i], true
	}
	return ObjectNodeKey{}, false
}
