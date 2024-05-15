package constraint

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
)

const colonTrue = ": true"
const colonFalse = ": false"

func parseUint(v bytes.Bytes, c Type) uint {
	u, err := v.ParseUint()
	if err != nil {
		panic(errs.ErrInvalidValueOfConstraint.F(c.String()))
	}
	return u
}
