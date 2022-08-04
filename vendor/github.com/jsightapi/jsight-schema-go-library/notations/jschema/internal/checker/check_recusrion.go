package checker

import (
	"strings"

	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/notations/jschema/internal/schema"
)

// CheckRecursion checks that given schema doesn't have invalid recursions.
//
// Examples of invalid recursions:
// TYPE @foo
// {
//   "foo": @foo
// }
//
// Examples of valid recursions:
// TYPE @foo
// {
//   "foo": @foo // {optional: true}
// }
//
// TYPE @foo
// {
//   "foo": [@foo]
// }
func CheckRecursion(rootTypeName string, rootSchema *schema.Schema) error {
	if rootSchema.RootNode() == nil {
		return nil
	}

	rc := &recursionChecker{
		visited: map[string]struct{}{
			// Obviously, root type was visited.
			rootTypeName: {},
		},
		path: []string{rootTypeName},
	}

	return rc.check(rootSchema.RootNode(), rootSchema.TypesList())
}

type recursionChecker struct {
	// visited a set of visited types.
	visited map[string]struct{}

	// path a path to current type.
	// Necessary for building an error message 'cause user should understand where
	// recursion was found.
	path []string
}

func (c *recursionChecker) check(node schema.Node, types map[string]schema.Type) error {
	// We can represent checked type and all dependent types as a graph. Root node
	// is a checked type. This is a directed graph with two types of edges: required
	// and optional.
	//
	// Edge between two types will be required if first type requires second type,
	// for instance when we have something like this:
	// {
	//   "foo": @bar
	// }
	// Here main type require `@bar` type.
	//
	// Edge between two types will be optional if first type doesn't require second
	// type, for instance:
	// {
	//   "foo": @bar // {optional: true}
	// }
	// Here main type doesn't require `@bar` type.
	//
	// So, we will visit all required types and skip optional and mark all passed
	// types until we pass all types. Or we try to mark already marked type. In
	// that case we have infinity recursion.
	//
	// Example:
	//
	// TYPE @foo
	// {
	//   "foo1": [@foo],
	//   "foo2": @bar
	// }
	//
	// TYPE @bar
	// {
	//   "bar1": @fizz
	//   "bar2": @bar // {optional: true}
	// }
	//
	// TYPE @fizz
	// {
	//   "fizz1": @foo
	// }
	//
	// Here we check `@foo` type. We can omit all optional dependencies, and we
	// will get next situation:
	//
	// @foo -> @bar -> @fizz -> @foo -> @bar -> @fizz -> @foo ...
	//
	// Obviously, this is an infinity recursion. And no mater which type we will
	// take to check we will get an infinity recursion here.

	// Ignore optional node.
	if schema.IsOptionalNode(node) {
		return nil
	}

	switch node := node.(type) {
	// Array might contain no items, so it's optional and should be skipped.
	// Literal nodes should be skipped 'cause it doesn't contain any fields or
	// type names.
	// Mixed node doesn't contain user type.
	case *schema.ArrayNode, *schema.LiteralNode, *schema.MixedNode:
		return nil

	// Special logic for mixed value 'cause it can contain a link to another type.
	case *schema.MixedValueNode:
		return c.checkMixedValueNode(node, types)

	// We should check all fields in the object 'cause some of them can be required.
	case *schema.ObjectNode:
		for _, n := range node.Children() {
			if err := c.check(n, types); err != nil {
				return err
			}
		}

	default:
		return errors.ErrImpossible
	}

	return nil
}

func (c *recursionChecker) checkMixedValueNode(
	node *schema.MixedValueNode,
	types map[string]schema.Type,
) error {
	tt := node.GetTypes()

	// We should check all types and return an error only if all paths leads
	// to infinity recursion.
	errs := make([]error, 0, len(tt))
	for _, t := range tt {
		if err := c.checkType(t, types); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 && len(errs) == len(tt) {
		// Just return first found error.
		return errs[0]
	}
	return nil
}

func (c *recursionChecker) checkType(typeName string, types map[string]schema.Type) error {
	if !c.visit(typeName) {
		return c.createError()
	}
	defer c.leave(typeName)

	t := types[typeName]
	if t.Schema() == nil {
		// This might happen if we didn't know anything about this type.
		// Normally we shouldn't get this situation.
		return nil
	}

	return c.check(t.Schema().RootNode(), t.Schema().TypesList())
}

func (c *recursionChecker) visit(typeName string) bool {
	c.path = append(c.path, typeName)
	if _, ok := c.visited[typeName]; ok {
		return false
	}
	c.visited[typeName] = struct{}{}
	return true
}

func (c *recursionChecker) leave(typeName string) {
	if len(c.path) > 0 {
		c.path = c.path[:len(c.path)-1]
	}
	delete(c.visited, typeName)
}

func (c *recursionChecker) createError() error {
	return errors.Format(errors.ErrInfinityRecursionDetected, strings.Join(c.path, " -> "))
}
