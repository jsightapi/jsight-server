package constraint

import (
	stdBytes "bytes"

	schema "github.com/jsightapi/jsight-schema-core"
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
	"github.com/jsightapi/jsight-schema-core/json"
)

type UUID struct{}

var (
	_ Constraint       = UUID{}
	_ Constraint       = (*UUID)(nil)
	_ LiteralValidator = UUID{}
	_ LiteralValidator = (*UUID)(nil)
)

func NewUuid() *UUID {
	return &UUID{}
}

func (UUID) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (UUID) Type() Type {
	return UuidConstraintType
}

func (UUID) String() string {
	return UuidConstraintType.String()
}

func (UUID) Validate(value bytes.Bytes) {
	err := parseBytes(value.Unquote().Data())
	if err != nil {
		panic(errs.ErrInvalidUUID.F(err))
	}
}

func (UUID) ASTNode() schema.RuleASTNode {
	return newEmptyRuleASTNode()
}

// parseBytes parse UUID bytes.
// All the following on the basis of github.com/google/uuid.
func parseBytes(b []byte) error { //nolint:gocyclo // For now it's okay.
	switch len(b) {
	case 36: // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9: // urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		if !stdBytes.Equal(stdBytes.ToLower(b[:9]), []byte("urn:uuid:")) {
			return errs.ErrURNPrefix.F(b[:9])
		}
		b = b[9:]
	case 36 + 2: // {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
		if b[0] != '{' || b[37] != '}' {
			return errs.ErrUUIDPrefix.F()
		}
		b = b[1:]
	case 32: // xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		for i := 0; i < 32; i += 2 {
			if !xtob(b[i], b[i+1]) {
				return errs.ErrUUIDFormat.F()
			}
		}
		return nil
	default:
		return errs.ErrUUIDLength.F(len(b))
	}

	// it must be of the form  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return errs.ErrUUIDFormat.F()
	}
	for _, x := range [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		if !xtob(b[x], b[x+1]) {
			return errs.ErrUUIDFormat.F()
		}
	}
	return nil
}

func xtob(x1, x2 byte) bool {
	b1 := xvalues[x1]
	b2 := xvalues[x2]
	return b1 != 255 && b2 != 255
}

var xvalues = [256]byte{
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 10, 11, 12, 13, 14, 15, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
	255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
}
