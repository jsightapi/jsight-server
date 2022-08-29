package validator

import (
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/lexeme"
)

// A Tree of validators.
//
// The tree consists of nodes (validators).
// Each node contains a pointer to the parent and validator from the schema package.
// Leaf, is the node without children.
// Root, is the top node in a tree. The parent of the root node is nil
//
// The "tree" structure contains the leaves in which to pass the LexEvent for
// validation.
type Tree struct {
	// leaves a list of all the leaves of the tree.
	leaves map[int]validator

	// leavesIndexes the list of indexes for the leaves. Defined in struct (not
	// in method) to optimize memory allocation.
	leavesIndexes []int

	// nextIndex the index for next leaf.
	nextIndex int
}

func NewTree(list []validator) Tree {
	t := Tree{
		nextIndex:     0,
		leaves:        make(map[int]validator, 5),
		leavesIndexes: make([]int, 0, 5),
	}
	for _, v := range list {
		t.addLeaf(v)
	}

	return t
}

// FeedLeaves returns true if the validation of the entire tree is completed. There
// are no more validators left in the tree.
func (t *Tree) FeedLeaves(jsonLex lexeme.LexEvent) bool {
	// A new array with the indexes of the leaves of the tree, to iterate on it.
	// The tree will change during the iteration.
	t.setLeavesIndexes()
	errorsCount := 0

	var err error

	for _, indexOfLeaf := range t.leavesIndexes {
		if leaf, ok := t.leaves[indexOfLeaf]; ok {
			err = t.feedLeaf(leaf, jsonLex, indexOfLeaf) // can panic
			if err != nil {
				errorsCount++
			}
		}
	}

	if errorsCount == len(t.leavesIndexes) {
		if len(t.leavesIndexes) == 1 {
			panic(err)
		} else {
			panic(lexeme.NewLexEventError(jsonLex, errors.ErrOrRuleSetValidation))
		}
	}

	if len(t.leaves) == 0 {
		return true
	}
	return false
}

func (t *Tree) setLeavesIndexes() {
	t.leavesIndexes = t.leavesIndexes[:0]
	for i := range t.leaves {
		t.leavesIndexes = append(t.leavesIndexes, i)
	}
}

// feedLeaf passes the LexEvent to the validator. Based on the results changes the
// tree.
// Removes or adds new validators to the tree.
// Returns common.DocumentError if an error is found during node validation.
func (t *Tree) feedLeaf(leaf validator, jsonLex lexeme.LexEvent, indexOfLeaf int) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(errors.DocumentError)
			if ok {
				delete(t.leaves, indexOfLeaf)
			} else {
				panic(r)
			}
		}
	}()

	children, done := leaf.feed(jsonLex) // can panic

	if done { // validation of node completed
		parent := leaf.parent()
		leaf.setParent(nil) // remove the pointer to simplify garbage collection in the future
		if parent == nil {
			delete(t.leaves, indexOfLeaf)
		} else {
			t.leaves[indexOfLeaf] = parent // step back to parent
		}
		return nil
	}

	// children found
	for j, child := range children {
		if j == 0 {
			// Forget/replace the current leaf. He becomes branch, parent for
			// first child.
			t.leaves[indexOfLeaf] = child
		} else {
			t.addLeaf(child) // append new child leaf to tree
		}
	}

	return nil
}

func (t *Tree) addLeaf(v validator) {
	t.leaves[t.nextIndex] = v
	t.nextIndex++
}
