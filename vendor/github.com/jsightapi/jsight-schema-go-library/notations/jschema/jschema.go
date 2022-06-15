package jschema

import (
	stdErrors "errors"
	"fmt"
	"io"
	"strings"
	"sync"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/formats/json"
	"github.com/jsightapi/jsight-schema-go-library/fs"
	"github.com/jsightapi/jsight-schema-go-library/notations/internal"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/checker"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/loader"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/scanner"
	internalSchema "github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema/constraint"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/validator"
	"github.com/jsightapi/jsight-schema-go-library/notations/regex"
)

type Schema struct {
	file  *fs.File
	inner *internalSchema.Schema

	loadErr    error
	compileErr error

	usedUserTypes []string

	astNode jschema.ASTNode

	len uint

	lenOnce     sync.Once
	loadOnce    sync.Once
	compileOnce sync.Once

	allowTrailingNonSpaceCharacters bool
	areKeysOptionalByDefault        bool
}

var _ jschema.Schema = &Schema{}

// New creates a Jsight schema with specified name and content.
func New(name string, content []byte, oo ...Option) *Schema {
	return FromFile(fs.NewFile(name, content), oo...)
}

// FromFile creates a Jsight schema from file.
func FromFile(f *fs.File, oo ...Option) *Schema {
	s := &Schema{
		file: f,
	}

	for _, o := range oo {
		o(s)
	}

	return s
}

type Option func(s *Schema)

func AllowTrailingNonSpaceCharacters() Option {
	return func(s *Schema) {
		s.allowTrailingNonSpaceCharacters = true
	}
}

func KeysAreOptionalByDefault() Option {
	return func(s *Schema) {
		s.areKeysOptionalByDefault = true
	}
}

func (s *Schema) Len() (uint, error) {
	var err error
	s.lenOnce.Do(func() {
		s.len, err = s.computeLen()
	})
	return s.len, err
}

func (s *Schema) computeLen() (length uint, err error) {
	// Iterate through all lexemes until we reach the end
	// We should rewind here in case we call NextLexeme method.
	defer func() {
		err = handlePanic(recover(), err)
	}()

	return scanner.NewSchemaScanner(s.file, true).Length(), err
}

func (s *Schema) Example() (b []byte, err error) {
	defer func() {
		err = handlePanic(recover(), err)
	}()

	if err := s.load(); err != nil {
		return nil, err
	}

	if s.inner.RootNode() == nil {
		return nil, errors.NewDocumentError(s.file, errors.ErrEmptySchema)
	}

	return buildExample(s.inner.RootNode()), nil
}

func buildExample(node internalSchema.Node) []byte {
	c := node.Constraint(constraint.TypesListConstraintType)
	if c != nil {
		panic(errors.ErrUserTypeFound)
	}

	switch typedNode := node.(type) {
	case *internalSchema.ObjectNode:
		b := make([]byte, 0, 512)
		b = append(b, '{')
		objectNode := node.(*internalSchema.ObjectNode) //nolint:errcheck // We're sure about this type.
		children := objectNode.Children()
		length := len(children)
		for i, childNode := range children {
			key := objectNode.Key(i)
			b = append(b, '"')
			b = append(b, []byte(key.Name)...)
			b = append(b, '"', ':')
			b = append(b, buildExample(childNode)...)
			if i+1 != length {
				b = append(b, ',')
			}
		}
		b = append(b, '}')
		return b

	case *internalSchema.ArrayNode:
		b := make([]byte, 0, 512)
		b = append(b, '[')
		children := typedNode.Children()
		length := len(children)
		for i, childNode := range children {
			b = append(b, buildExample(childNode)...)
			if i+1 != length {
				b = append(b, ',')
			}
		}
		b = append(b, ']')
		return b

	case *internalSchema.LiteralNode:
		return typedNode.BasisLexEventOfSchemaForNode().Value()

	default:
		panic(fmt.Sprintf("unhandled node type %T", node))
	}
}

func (s *Schema) AddType(name string, sc jschema.Schema) (err error) {
	defer func() {
		err = handlePanic(recover(), err)
	}()

	if err := s.load(); err != nil {
		return err
	}

	switch typ := sc.(type) {
	case *Schema:
		if err := typ.load(); err != nil {
			return fmt.Errorf("load added type: %w", err)
		}

		s.inner.AddNamedType(name, typ.inner, s.file, 0)
	case *regex.Schema:
		pattern, err := typ.Pattern()
		if err != nil {
			return err
		}

		example, err := typ.Example()
		if err != nil {
			return fmt.Errorf("generate example for Regex type: %w", err)
		}

		typSc := New(name, []byte(fmt.Sprintf("%q // {regex: %q}", example, pattern)))
		if err := typSc.load(); err != nil {
			return fmt.Errorf("load added type: %w", err)
		}

		s.inner.AddNamedType(name, typSc.inner, s.file, 0)

	default:
		return fmt.Errorf("schema should be JSight or Regex schema, but %T given", sc)
	}

	return nil
}

func (*Schema) AddRule(string, jschema.Rule) error {
	return stdErrors.New("not supported yet")
}

func (s *Schema) Check() (err error) {
	defer func() {
		err = handlePanic(recover(), err)
	}()
	return s.compile()
}

func (s *Schema) Validate(document jschema.Document) (err error) {
	defer func() {
		err = handlePanic(recover(), err)
	}()
	if err := s.compile(); err != nil {
		return err
	}

	if _, ok := document.(*json.Document); !ok {
		return fmt.Errorf("support only JSON documents, but got %T", document)
	}

	return s.validate(document)
}

func (s *Schema) validate(document jschema.Document) error {
	tree := validator.NewTree(
		validator.NodeValidatorList(s.inner.RootNode(), *s.inner, nil),
	)

	empty := true

	for {
		jsonLex, err := document.NextLexeme()
		if err != nil {
			if stdErrors.Is(err, io.EOF) {
				break
			}
			return err
		}

		empty = false
		if tree.FeedLeaves(jsonLex) { // can panic: error of validation
			break
		}
	}

	if empty {
		return internal.NewValidatorError(errors.ErrEmptyJson, "")
	}

	// check for error: Invalid non-space byte after top-level value
	for {
		_, err := document.NextLexeme()
		if err != nil {
			if stdErrors.Is(err, io.EOF) {
				break
			}
			return err
		}
	}
	return nil
}

func (s *Schema) GetAST() (an jschema.ASTNode, err error) {
	if err := s.compile(); err != nil {
		return jschema.ASTNode{}, err
	}

	return s.astNode, nil
}

func (s *Schema) UsedUserTypes() ([]string, error) {
	if err := s.load(); err != nil {
		return nil, err
	}
	return s.usedUserTypes, nil
}

func (s *Schema) load() error {
	s.loadOnce.Do(func() {
		defer func() {
			s.loadErr = handlePanic(recover(), s.loadErr)
		}()
		sc := loader.LoadSchemaWithoutCompile(
			scanner.NewSchemaScanner(s.file, s.allowTrailingNonSpaceCharacters),
			nil,
		)
		s.inner = &sc
		s.astNode = s.buildASTNode()
		s.collectUserTypes()
		loader.CompileBasic(s.inner, s.areKeysOptionalByDefault)
	})
	return s.loadErr
}

func (s *Schema) Build() error {
	return s.compile()
}

func (s *Schema) collectUserTypes() {
	node := s.inner.RootNode()
	// This is possible when schema isn't valid.
	if node == nil {
		return
	}

	uu := map[string]struct{}{}
	collectUserTypes(node, uu)
	s.usedUserTypes = make([]string, 0, len(uu))

	for u := range uu {
		s.usedUserTypes = append(s.usedUserTypes, u)
	}
}

func collectUserTypes(node internalSchema.Node, uu map[string]struct{}) {
	collectUserTypesFromTypesListConstraint(node, uu)
	collectUserTypesFromTypeConstraint(node, uu)
	collectUserTypesFromAllOfConstraint(node, uu)

	switch n := node.(type) {
	case *internalSchema.ObjectNode:
		collectUserTypesObjectNode(n, uu)

	case *internalSchema.ArrayNode:
		for _, child := range n.Children() {
			collectUserTypes(child, uu)
		}

	case *internalSchema.MixedValueNode:
		for _, ut := range strings.Split(n.Value().String(), "|") {
			s := strings.TrimSpace(ut)
			if s[0] == '@' {
				uu[s] = struct{}{}
			}
		}
	}
}

func collectUserTypesFromTypesListConstraint(node internalSchema.Node, uu map[string]struct{}) {
	if c := node.Constraint(constraint.TypesListConstraintType); c != nil {
		for _, name := range c.(*constraint.TypesList).Names() {
			if name[0] == '@' {
				uu[name] = struct{}{}
			}
		}
	}
}

func collectUserTypesFromTypeConstraint(node internalSchema.Node, uu map[string]struct{}) {
	if c := node.Constraint(constraint.TypeConstraintType); c != nil {
		s := c.(*constraint.TypeConstraint).Bytes().Unquote().String()
		if s[0] == '@' {
			uu[s] = struct{}{}
		}
	}
}

func collectUserTypesFromAllOfConstraint(node internalSchema.Node, uu map[string]struct{}) {
	if c := node.Constraint(constraint.AllOfConstraintType); c != nil {
		for _, s := range c.(*constraint.AllOf).SchemaNames() {
			if s[0] == '@' {
				uu[s] = struct{}{}
			}
		}
	}
}

func collectUserTypesObjectNode(node *internalSchema.ObjectNode, uu map[string]struct{}) {
	node.Keys().EachSafe(func(k string, v internalSchema.InnerObjectNodeKey) {
		if v.IsShortcut {
			if k[0] == '@' {
				uu[k] = struct{}{}
			}
		}

		c, ok := node.Child(k)
		if ok {
			collectUserTypes(c, uu)
		}
	})
}

func (s *Schema) buildASTNode() jschema.ASTNode {
	root := s.inner.RootNode()
	if root == nil {
		// This case will be handled in loader.CompileBasic.
		return jschema.ASTNode{
			Properties: &jschema.ASTNodes{},
			Rules:      &jschema.RuleASTNodes{},
		}
	}

	an, err := root.ASTNode()
	if err != nil {
		panic(err)
	}
	return an
}

func (s *Schema) compile() error {
	s.compileOnce.Do(func() {
		defer func() {
			s.compileErr = handlePanic(recover(), s.compileErr)
		}()
		if err := s.load(); err != nil {
			s.compileErr = err
			return
		}
		loader.CompileAllOf(s.inner)
		loader.AddUnnamedTypes(s.inner)
		checker.CheckRootSchema(s.inner)
	})
	return s.compileErr
}

func handlePanic(r interface{}, originErr error) error {
	if originErr != nil {
		return originErr
	}

	if r == nil {
		return nil
	}

	rErr, ok := r.(error)
	if !ok {
		panic(r)
	}
	return rErr
}
