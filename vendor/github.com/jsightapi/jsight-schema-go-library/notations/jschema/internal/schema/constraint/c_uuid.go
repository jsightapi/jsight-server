package constraint

import (
	"bytes"
	stdErrors "errors"
	"fmt"

	jschema "github.com/jsightapi/jsight-schema-go-library"
	jbytes "github.com/jsightapi/jsight-schema-go-library/bytes"
	"github.com/jsightapi/jsight-schema-go-library/errors"
	"github.com/jsightapi/jsight-schema-go-library/internal/json"
)

type Uuid struct{}

var _ Constraint = Uuid{}

func NewUuid() *Uuid {
	return &Uuid{}
}

func (Uuid) IsJsonTypeCompatible(t json.Type) bool {
	return t == json.TypeString
}

func (Uuid) Type() Type {
	return UuidConstraintType
}

func (Uuid) String() string {
	return UuidConstraintType.String()
}

func (Uuid) Validate(value jbytes.Bytes) {
	err := ParseBytes(value.Unquote())
	if err != nil {
		panic(errors.Format(errors.ErrInvalidUuid, err))
	}
}

func (Uuid) ASTNode() jschema.RuleASTNode {
	return newEmptyRuleASTNode()
}

// all of the following on the basis of github.com/google/uuid

func ParseBytes(b []byte) error { //nolint:gocyclo // For now it's okay.
	switch len(b) {
	case 36: // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9: // urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		if !bytes.Equal(bytes.ToLower(b[:9]), []byte("urn:uuid:")) {
			return fmt.Errorf("invalid urn prefix: %q", b[:9])
		}
		b = b[9:]
	case 36 + 2: // {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
		if b[0] != '{' || b[37] != '}' {
			return stdErrors.New("invalid prefix: braces expected")
		}
		b = b[1:]
	case 32: // xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		for i := 0; i < 32; i += 2 {
			if !xtob(b[i], b[i+1]) {
				return stdErrors.New("invalid UUID format")
			}
		}
		return nil
	default:
		return fmt.Errorf("invalid UUID length: %d", len(b))
	}

	// it must be of the form  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' {
		return stdErrors.New("invalid UUID format")
	}
	for _, x := range [16]int{
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34} {
		if !xtob(b[x], b[x+1]) {
			return stdErrors.New("invalid UUID format")
		}
	}
	return nil
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

func xtob(x1, x2 byte) bool {
	b1 := xvalues[x1]
	b2 := xvalues[x2]
	return b1 != 255 && b2 != 255
}
