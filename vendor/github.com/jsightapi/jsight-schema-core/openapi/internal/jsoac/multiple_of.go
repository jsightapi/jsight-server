package jsoac

import (
	schema "github.com/jsightapi/jsight-schema-core"

	"math"

	"strconv"
)

func newMultipleOf(astNode schema.ASTNode) *float64 {
	if astNode.Rules.Has("precision") {
		v := astNode.Rules.GetValue("precision").Value
		precision, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil
		}

		// from api schema validation, precision can not be zero
		multipleOf := 1 / math.Pow(10, precision)
		if multipleOf == 0 {
			return nil
		}

		return &multipleOf
	}
	return nil
}
