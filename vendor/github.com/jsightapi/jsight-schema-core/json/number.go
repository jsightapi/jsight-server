package json

import (
	"strconv"

	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
)

// Number provides a method to work with numbers in the understanding of JSON (for
// example -1.2 e+3).
type Number struct {
	nat bytes.Bytes

	// exp an absolute value of exponent. In fact, a negative value is assumed,
	// but for optimization it is stored without a minus sign.
	exp int

	neg bool
}

func NewNumber(b bytes.Bytes) (*Number, error) {
	return newScanner().Scan(b)
}

// trimLeadingZerosInTheIntegerPart removes zeros from the beginning of the integer
// part (if any).
func (n *Number) trimLeadingZerosInTheIntegerPart() error {
	length := n.nat.Len()
	if n.exp < 0 || n.exp > length {
		return errs.ErrIncorrectExponentValue.F()
	}
	for intLen := length - n.exp; intLen != 0; intLen-- {
		c := n.nat.FirstByte()
		if c != '0' {
			break
		}
		n.nat = n.nat.SubLow(1) // trim left byte
	}
	return nil
}

// trimTrailingZerosInTheFractionalPart removes zeros from the end of the fractional
// part (if any).
func (n *Number) trimTrailingZerosInTheFractionalPart() error {
	if n.exp < 0 || n.exp > n.nat.Len() {
		return errs.ErrIncorrectExponentValue.F()
	}
	for ; n.exp != 0; n.exp-- {
		i := n.nat.Len() - 1
		c := n.nat.Byte(i) // last character
		if c != '0' {
			break
		}
		n.nat = n.nat.SubHigh(bytes.Index(i)) // trim right byte
	}
	return nil
}

func (n Number) int() bytes.Bytes {
	return n.nat.SubHigh(bytes.Index(n.nat.Len() - n.exp))
}

func (n Number) fra() bytes.Bytes {
	return n.nat.SubLow(bytes.Index(n.nat.Len() - n.exp))
}

func (n Number) LengthOfFractionalPart() uint {
	return uint(n.exp)
}

// Cmp compares the numbers represented by n and nn and returns:
//
//	-1 if n <  nn
//	 0 if n == nn
//	+1 if n >  nn.
func (n Number) Cmp(nn *Number) int {
	if n.neg == nn.neg {
		b := n.cmpAbs(nn)
		if n.neg {
			return n.not(b)
		}
		return b
	}

	if n.neg {
		return -1
	}
	return 1
}

func (Number) not(cmp int) int {
	switch cmp {
	case 1:
		return -1
	case -1:
		return 1
	case 0:
		return 0
	}
	panic(errs.ErrRuntimeFailure.F())
}

func (n Number) cmpAbs(nn *Number) (r int) {
	cmp := n.cmpInt(nn)
	if cmp == 0 {
		return n.cmpFra(nn)
	}
	return cmp
}

func (n Number) cmpInt(nn *Number) (r int) {
	x := n.int()
	y := nn.int()
	xLen := x.Len()
	yLen := y.Len()
	if xLen != yLen || xLen == 0 {
		switch {
		case xLen < yLen:
			r = -1
		case xLen > yLen:
			r = 1
		}
		return
	}

	// xLen == yLen
	for i := 0; i < xLen; i++ {
		switch {
		case x.Byte(i) < y.Byte(i):
			return -1
		case x.Byte(i) > y.Byte(i):
			return 1
		}
	}

	return 0
}

func (n Number) cmpFra(nn *Number) (r int) {
	x := n.fra()
	y := nn.fra()
	xLen := x.Len()
	yLen := y.Len()

	var length int
	if xLen > yLen {
		length = xLen
	} else {
		length = yLen
	}

	for i := 0; i < length; i++ {
		digit1 := 0
		digit2 := 0
		if i < xLen {
			digit1 = int(x.Byte(i)) - 48
		}
		if i < yLen {
			digit2 = int(y.Byte(i)) - 48
		}
		if digit1 < digit2 {
			return -1
		} else if digit1 > digit2 {
			return 1
		}
	}

	return 0
}

// Equal returns whether the numbers represented by n and nn are equal.
func (n Number) Equal(nn *Number) bool {
	return n.Cmp(nn) == 0
}

// GreaterThan (GT) returns true when n is greater than nn.
func (n Number) GreaterThan(nn *Number) bool {
	return n.Cmp(nn) == 1
}

// GreaterThanOrEqual (GTE) returns true when n is greater than or equal to nn.
func (n Number) GreaterThanOrEqual(nn *Number) bool {
	cmp := n.Cmp(nn)
	return cmp == 1 || cmp == 0
}

// LessThan (LT) returns true when n is less than nn.
func (n Number) LessThan(nn *Number) bool {
	return n.Cmp(nn) == -1
}

// LessThanOrEqual (LTE) returns true when n is less than or equal to nn.
func (n Number) LessThanOrEqual(nn *Number) bool {
	cmp := n.Cmp(nn)
	return cmp == -1 || cmp == 0
}

func (n Number) String() string {
	var str string
	i := n.int()
	f := n.fra()
	if n.neg {
		str += "-"
	}
	if i.Len() == 0 {
		str += "0"
	} else {
		str += i.String()
	}
	if f.Len() != 0 {
		str += "." + f.String()
	}
	return str
}

func (n Number) ToFloat() float64 {
	v, err := strconv.ParseFloat(n.String(), 64)
	// Normally we shouldn't get an error here cause string value is valid.
	// But we should throw error just in case.
	if err != nil {
		panic(errs.ErrRuntimeFailure.F())
	}
	return v
}
