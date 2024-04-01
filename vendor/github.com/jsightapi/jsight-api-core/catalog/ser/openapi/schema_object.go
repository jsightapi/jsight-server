package openapi

import (
	"github.com/jsightapi/jsight-api-core/catalog"
	"github.com/jsightapi/jsight-api-core/notation"
)

func schemaObjectFromExchangeSchema(es catalog.ExchangeSchema) schemaObject {
	switch es.Notation() {
	case notation.SchemaNotationJSight, notation.SchemaNotationRegex:
		return schemaObjectFromSchema(es)
	case notation.SchemaNotationAny:
		return schemaObjectForAny()
	case notation.SchemaNotationEmpty:
		panic("notation 'empty' cannot be represented by SchemaObject")
	default:
		panic("unsupported schema notation")
	}
}

func schemaObjectFromSchema(es catalog.ExchangeSchema) schemaObject {
	switch s := es.(type) {
	case *catalog.ExchangeJSightSchema:
		return getJSchemaInfo(s.JSchema).SchemaObject()
	case *catalog.ExchangeRegexSchema:
		return getRSchemaInfo(s.RSchema).SchemaObject()
	default:
		panic("unsupported ExchangeSchema type for this method")
	}
}

func schemaObjectForAny() schemaObject {
	return &schemaObjectAny{}
}
