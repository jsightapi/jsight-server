package catalog

import (
	schema "github.com/jsightapi/jsight-schema-core"

	"github.com/jsightapi/jsight-api-core/directive"
)

type Prop struct {
	Parameter string
	ASTNode   schema.ASTNode
	Directive directive.Directive
}
