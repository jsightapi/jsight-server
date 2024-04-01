package openapi

import (
	"encoding/json"
	"fmt"

	"github.com/jsightapi/jsight-api-core/catalog"
)

//nolint:unused
func debugExchangeSchema(es catalog.ExchangeSchema) {
	fmt.Printf("es type is %T\n", es)
	fmt.Printf("es notation is %s\n", es.Notation())

	switch s := es.(type) {
	case *catalog.ExchangeJSightSchema:
		ast, err := s.JSchema.GetAST()
		if err != nil {
			fmt.Printf("No AST :(\n")
		} else {
			fmt.Printf("Value: %s \n", ast.Value)
		}
	default:
		fmt.Printf("Not a JSchema\n")
	}
	fmt.Printf("Debug Over\n")
}

//nolint:unused
func schemaObjectToString(so schemaObject) string {
	b, err := json.Marshal(so)
	if err == nil {
		return string(b)
	} else {
		return "marshaling error"
	}
}
