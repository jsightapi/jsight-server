package ischema

import (
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/lexeme"
)

type ObjectNodeKeys struct {
	index map[indexKey]int
	Data  []ObjectNodeKey
}

type ObjectNodeKey struct {
	Key        string
	Lex        lexeme.LexEvent
	Index      int
	IsShortcut bool
}

type indexKey struct {
	Key        string
	IsShortcut bool
}

func indexKeyFromObjectNodeKey(k ObjectNodeKey) indexKey {
	return indexKey{
		Key:        k.Key,
		IsShortcut: k.IsShortcut,
	}
}

func newObjectNodeKeys() *ObjectNodeKeys {
	return &ObjectNodeKeys{
		Data:  make([]ObjectNodeKey, 0, 5),
		index: make(map[indexKey]int, 5),
	}
}

func (k *ObjectNodeKeys) Set(v ObjectNodeKey) {
	if k.isDuplicatedKey(v) {
		panic(errs.ErrDuplicateKeysInSchema.F(v.Key))
	}

	k.index[indexKeyFromObjectNodeKey(v)] = v.Index
	k.Data = append(k.Data, v)
}

func (k *ObjectNodeKeys) isDuplicatedKey(newKey ObjectNodeKey) bool {
	_, ok := k.index[indexKeyFromObjectNodeKey(newKey)]
	return ok
}

func (k ObjectNodeKeys) Find(i int) (ObjectNodeKey, bool) {
	if len(k.Data) > i {
		return k.Data[i], true
	}
	return ObjectNodeKey{}, false
}

func (k ObjectNodeKeys) Get(key string, isShortcut bool) (ObjectNodeKey, bool) {
	if i, ok := k.index[indexKey{
		Key:        key,
		IsShortcut: isShortcut,
	}]; ok {
		return k.Data[i], true
	}
	return ObjectNodeKey{}, false
}
