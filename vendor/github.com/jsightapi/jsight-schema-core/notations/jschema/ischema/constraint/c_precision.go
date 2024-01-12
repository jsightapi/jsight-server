package constraint

import (
	"strconv"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type Precision struct {
	value uint
}

var (
	_ Constraint       = Precision{}
	_ Constraint       = (*Precision)(nil)
	_ LiteralValidator = Precision{}
	_ LiteralValidator = (*Precision)(nil)
)

func NewPrecision(ruleValue bytes.Bytes) *Precision {
	u := parseUint(ruleValue, PrecisionConstraintType)

	if u == 0 {
		panic(errs.ErrZeroPrecision.F())
	}

	return &Precision{
		value: u,
	}
}

func (Precision) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeFloat
}

func (Precision) Type() Type {
	return PrecisionConstraintType
}

func (c Precision) String() string {
	return PrecisionConstraintType.String() + ": " + strconv.Itoa(int(c.value))
}

func (c Precision) Validate(value bytes.Bytes) {
	n, err := json.NewNumber(value)
	if err != nil {
		panic(err)
	}
	if c.value < n.LengthOfFractionalPart() {
		panic(errs.ErrConstraintValidation.F(
			PrecisionConstraintType.String(),
			strconv.Itoa(int(c.value)),
			"(exclusive)",
		))
	}
}

func (c Precision) ASTNode() schema.RuleASTNode {
	return newRuleASTNode(
		schema.TokenTypeNumber,
		strconv.FormatUint(uint64(c.value), 10),
		schema.RuleASTNodeSourceManual,
	)
}
