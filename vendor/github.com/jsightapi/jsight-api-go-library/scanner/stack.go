package scanner

import (
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/fs"

	"github.com/jsightapi/jsight-api-go-library/directive"
	"github.com/jsightapi/jsight-api-go-library/jerr"
)

// Stack ths stack of scanners.
type Stack struct {
	uniqueFiles    map[string]struct{}
	includeTracers map[uint64]directive.IncludeTracer
	stack          []stackItem
	hashes         []uint64
}

var ErrRecursionDetected = errors.New("recursion detected")

type stackItem struct {
	// The scanner which was pushed into stack.
	scanner *Scanner

	// at represents a position in on which we stop in the pushed scanner.
	at bytes.Index
}

// Push pushed new scanner into stack with position in this scanner where it was
// pushed. This position will be necessary for creating stack trace.
func (s *Stack) Push(scanner *Scanner, at bytes.Index) error {
	if s.uniqueFiles == nil {
		s.uniqueFiles = map[string]struct{}{}
	}

	name := scanner.file.Name()

	if _, ok := s.uniqueFiles[name]; ok {
		return ErrRecursionDetected
	}

	s.stack = append(s.stack, stackItem{
		scanner: scanner,
		at:      at,
	})
	s.uniqueFiles[name] = struct{}{}

	hash, err := s.computeScannerHash(scanner)
	if err != nil {
		return fmt.Errorf("failed to compute scanner's hash: %w", err)
	}
	s.hashes = append(s.hashes, hash)

	return nil
}

func (*Stack) computeScannerHash(scanner *Scanner) (uint64, error) {
	hasher := fnv.New64()
	path := scanner.File().Name()

	n, err := hasher.Write([]byte(path))
	if err != nil {
		return 0, err
	}

	if n != len(path) {
		return 0, errors.New("written less bytes then expected")
	}

	return hasher.Sum64(), nil
}

// Pop pops last add scanner.
func (s *Stack) Pop() *Scanner {
	l := len(s.stack)
	if l == 0 {
		return nil
	}
	e := s.stack[l-1].scanner
	s.stack = s.stack[:l-1]
	delete(s.uniqueFiles, e.file.Name())
	s.hashes = s.hashes[:l-1]
	return e
}

// Empty returns true is stack is empty.
func (s *Stack) Empty() bool {
	return len(s.stack) == 0
}

// AddIncludeTraceToError adds include trace to given error.
func (s *Stack) AddIncludeTraceToError(je *jerr.JApiError) {
	addIncludeTraceToError(je, s.stack)
}

// ToDirectiveIncludeTracer converts this scanner stack to directive's include trace.
// We can't pass this stack itself 'cause it will be empty at the time when it
// will be necessary.
func (s *Stack) ToDirectiveIncludeTracer() directive.IncludeTracer {
	// Include tracer won't be used in normal situation by really necessary when
	// everything goes wrong. But it will impact on memory footprints. So we should
	// think about performance here.

	// Fast path.
	if s.Empty() {
		return emptyTracer
	}

	if s.includeTracers == nil {
		s.includeTracers = map[uint64]directive.IncludeTracer{}
	}

	currHash := s.hashes[len(s.hashes)-1]
	if tracer, ok := s.includeTracers[currHash]; ok {
		return tracer
	}

	// Slow path.
	tracer := newDirectiveIncludeTracer(s.stack)
	s.includeTracers[currHash] = tracer

	return tracer
}

var emptyTracer = emptyIncludeTracer{}

type emptyIncludeTracer struct{}

func (emptyIncludeTracer) AddIncludeTraceToError(*jerr.JApiError) {}

// directiveIncludeTracer represent an implementation of directive.IncludeTracer.
type directiveIncludeTracer struct {
	stack []directiveIncludeTracerItem
}

type directiveIncludeTracerItem struct {
	file *fs.File
	at   bytes.Index
}

func newDirectiveIncludeTracer(ii []stackItem) directiveIncludeTracer {
	d := directiveIncludeTracer{
		stack: make([]directiveIncludeTracerItem, 0, len(ii)),
	}
	for _, i := range ii {
		d.stack = append(d.stack, directiveIncludeTracerItem{
			file: i.scanner.file,
			at:   i.at,
		})
	}
	return d
}

func (d directiveIncludeTracer) AddIncludeTraceToError(je *jerr.JApiError) {
	addIncludeTraceToError(je, d.stack)
}

func addIncludeTraceToError[T any](je *jerr.JApiError, stack []T) {
	if je == nil {
		return
	}

	sl := len(stack)
	if sl == 0 {
		return
	}

	if je.HasStackTrace() {
		// We shouldn't add stack trace if error already have it.
		return
	}

	for i := sl - 1; i >= 0; i-- {
		switch e := any(stack[i]).(type) {
		case stackItem:
			je.OccurredInFile(e.scanner.file, e.at)
		case directiveIncludeTracerItem:
			je.OccurredInFile(e.file, e.at)
		}
	}
}
