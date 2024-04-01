package jsoac

import (
	"encoding/json"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/internal/sync"

	"strconv"
)

const stringTrue = "true"
const stringFalse = "false"
const stringNull = "null"
const stringAny = "any"
const stringEnum = "enum"
const stringArray = "array"

var bufferPool = sync.NewBufferPool(1024)

// toJSONString returns JSON quoted string data
func toJSONString(s string) []byte {
	b, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return b
}

func isNullable(astNode schema.ASTNode) bool {
	if astNode.Rules.Has("nullable") && astNode.Rules.GetValue("nullable").Value == stringTrue {
		return true
	}
	return false
}

func isString(astNode schema.ASTNode) bool {
	return astNode.TokenType == schema.TokenTypeString
}

func int64Ref(i int64) *int64 {
	return &i
}

func int64RefByString(s string) *int64 {
	value, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	return int64Ref(value)
}
