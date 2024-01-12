package schema

// Schema represents a schema.
// Schema is a some description of expected structure of payload.
type Schema interface {
	// Len returns length of this schema in bytes.
	// Might return ParsingError if schema isn't valid.
	Len() (uint, error)

	// Example returns an example for this schema.
	// Might return ParsingError if schema isn't valid.
	Example() ([]byte, error)

	// AddType adds a new type to this schema.
	// Might return a ParsingError if add type isn't valid.
	AddType(name string, schema Schema) error

	// AddRule adds a new type to this schema.
	// Might return a ParsingError if add type isn't valid.
	AddRule(name string, schema Rule) error

	// Check checks that this schema is valid.
	// Can return ParsingError if schema isn't valid.
	Check() error

	// GetAST returns a root AST node for this schema.
	GetAST() (ASTNode, error)

	// UsedUserTypes return all used user types.
	UsedUserTypes() ([]string, error)
}
