package json

import (
	"errors"
	"strconv"

	"github.com/jsightapi/jsight-schema-go-library/bytes"
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

func NewNumberFromUint(u uint) *Number {
	n := Number{
		neg: false,
		nat: bytes.Bytes(strconv.FormatUint(uint64(u), 10)),
		exp: 0,
	}
	if err := n.trimLeadingZerosInTheIntegerPart(); err != nil {
		panic(err)
	}
	return &n
}

func NewNumberFromInt(i int) *Number {
	var natural bytes.Bytes
	negative := false
	if i < 0 {
		negative = true
		natural = bytes.Bytes(strconv.FormatUint(uint64(-i), 10))
	} else {
		natural = bytes.Bytes(strconv.FormatUint(uint64(i), 10))
	}
	n := Number{
		neg: negative,
		nat: natural,
		exp: 0,
	}
	if err := n.trimLeadingZerosInTheIntegerPart(); err != nil {
		panic(err)
	}
	return &n
}

func NewIntegerNumber(bytes bytes.Bytes) (*Number, error) {
	g := Guess(bytes)
	if !g.IsInteger() {
		return nil, errors.New(`Incorrect value "` + bytes.Unquote().String() + `". Must be an integer.`)
	}
	n, err := g.Number()
	if err != nil {
		return nil, errors.New(`Incorrect value "` + bytes.Unquote().String() + `". Unable to get number.`)
	}
	return n, nil
}

func NewNumber(b bytes.Bytes) (*Number, error) { //nolint:gocyclo // todo try to make this more readable
	type numberParseState int

	const (
		stateOnSearchStart numberParseState = iota
		stateMinusFound
		stateFirstZeroFound
		stateIntegerNumberFound
		statePointFound
		stateFractionalNumberFound
		stateExpFound
		stateExpSignFound
		stateExpNumberFound
	)

	var (
		state    numberParseState
		negative bool
		intLen   int
		fraLen   int
		expBegin int
	)

	for i, c := range b {
		switch state {
		case stateOnSearchStart:
			if c == '-' { //nolint:gocritic // todo rewrite this logic to switch
				negative = true
				state = stateMinusFound
				continue
			} else if c == '0' {
				intLen++
				state = stateFirstZeroFound
				continue
			} else if '1' <= c && c <= '9' {
				intLen++
				state = stateIntegerNumberFound
				continue
			}
		case stateMinusFound:
			if c == '0' {
				intLen++
				state = stateFirstZeroFound
				continue
			} else if '1' <= c && c <= '9' {
				intLen++
				state = stateIntegerNumberFound
				continue
			}
		case stateFirstZeroFound:
			if c == '.' {
				state = statePointFound
				continue
			}
		case stateIntegerNumberFound:
			if '0' <= c && c <= '9' { //nolint:gocritic // todo rewrite this logic to switch
				intLen++
				continue
			} else if c == '.' {
				state = statePointFound
				continue
			} else if c == 'e' || c == 'E' {
				state = stateExpFound
				continue
			}
		case statePointFound:
			if '0' <= c && c <= '9' {
				fraLen++
				state = stateFractionalNumberFound
				continue
			}
		case stateFractionalNumberFound:
			if '0' <= c && c <= '9' {
				fraLen++
				continue
			} else if c == 'e' || c == 'E' {
				state = stateExpFound
				continue
			}
		case stateExpFound:
			if c == '+' { //nolint:gocritic // todo rewrite this logic to switch
				state = stateExpSignFound
				continue
			} else if c == '-' {
				if expBegin == 0 {
					expBegin = i
				}
				state = stateExpSignFound
				continue
			} else if '0' <= c && c <= '9' {
				if expBegin == 0 {
					expBegin = i
				}
				continue
			}
		case stateExpSignFound:
			if '0' <= c && c <= '9' {
				if expBegin == 0 {
					expBegin = i
				}
				state = stateExpNumberFound
				continue
			}
		case stateExpNumberFound:
			if '0' <= c && c <= '9' {
				continue
			}
		}
		return nil, errors.New(`Incorrect number value "` + b.String() + `"`)
	}

	if state == stateOnSearchStart || state == stateMinusFound {
		return nil, errors.New(`Incorrect number value "` + b.String() + `"`)
	}

	if expBegin != 0 { // the EXP value found?
		exp, err := b[expBegin:].ParseInt()
		if err != nil {
			return nil, err
		}
		// example with negative exp: 12.34E-1 = 1.234; exp = -1; intLen = 2 + (-1) = 1
		// example with positive exp: 12.34E+1 = 123.4; exp =  1; intLen = 2 + 1    = 3
		intLen += exp
		fraLen -= exp
	}

	var natural bytes.Bytes

	// example 1.2E-2 = .012
	if intLen < 0 { //nolint:gocritic // todo rewrite this logic to switch
		natural = make(bytes.Bytes, 0, fraLen)
		natural = appendZeros(natural, -intLen)
		natural = appendDigits(b, natural)
	} else if fraLen < 0 { // example 1.2E+2 = 120
		natural = make(bytes.Bytes, 0, intLen)
		natural = appendDigits(b, natural)
		natural = appendZeros(natural, -fraLen)
		fraLen = 0
	} else { // example 12.3E-1 = 1.23
		natural = make(bytes.Bytes, 0, intLen+fraLen)
		natural = appendDigits(b, natural)
	}

	n := Number{
		neg: negative,
		nat: natural,
		exp: fraLen,
	}

	var err error

	err = n.trimLeadingZerosInTheIntegerPart()
	if err != nil {
		return nil, err
	}

	err = n.trimTrailingZerosInTheFractionalPart()
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func appendZeros(to bytes.Bytes, n int) bytes.Bytes {
	for ; n > 0; n-- {
		to = append(to, '0')
	}
	return to
}

func appendDigits(from bytes.Bytes, to bytes.Bytes) bytes.Bytes {
	for _, c := range from {
		if c == '-' || c == '.' { //nolint:gocritic // todo rewrite this logic to switch
			continue
		} else if '0' <= c && c <= '9' {
			to = append(to, c)
		} else { // c == 'e' || c == 'E'
			break
		}
	}
	return to
}

// Removes zeros from the beginning of the integer part (if any)
func (n *Number) trimLeadingZerosInTheIntegerPart() error {
	length := len(n.nat)
	if n.exp < 0 || n.exp > length {
		return errors.New("incorrect exponent value")
	}
	for intLen := length - n.exp; intLen != 0; intLen-- {
		c := n.nat[0] // first character
		if c != '0' {
			break
		}
		n.nat = n.nat[1:] // trim left byte
	}
	return nil
}

// Removes zeros from the end of the fractional part (if any)
func (n *Number) trimTrailingZerosInTheFractionalPart() error {
	if n.exp < 0 || n.exp > len(n.nat) {
		return errors.New("incorrect exponent value")
	}
	for ; n.exp != 0; n.exp-- {
		i := len(n.nat) - 1
		c := n.nat[i] // last character
		if c != '0' {
			break
		}
		n.nat = n.nat[:i] // trim right byte
	}
	return nil
}

func (n Number) int() bytes.Bytes {
	return n.nat[:len(n.nat)-n.exp]
}

func (n Number) fra() bytes.Bytes {
	return n.nat[len(n.nat)-n.exp:]
}

func (n Number) LengthOfFractionalPart() uint {
	return uint(n.exp)
}

// Cmp compares the numbers represented by n and nn and returns:
//  -1 if n <  nn
//   0 if n == nn
//  +1 if n >  nn
//
func (n Number) Cmp(nn *Number) int {
	if n.neg == nn.neg { //nolint:gocritic // todo rewrite this logic to switch
		b := n.cmpAbs(nn)
		if n.neg {
			return n.not(b)
		}
		return b
	} else if n.neg {
		return -1
	} else {
		return 1
	}
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
	panic("Incorrect value")
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
	xLen := len(x)
	yLen := len(y)
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
		case x[i] < y[i]:
			return -1
		case x[i] > y[i]:
			return 1
		}
	}

	return 0
}

func (n Number) cmpFra(nn *Number) (r int) {
	x := n.fra()
	y := nn.fra()
	xLen := len(x)
	yLen := len(y)

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
			digit1 = int(x[i]) - 48
		}
		if i < yLen {
			digit2 = int(y[i]) - 48
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
	if len(i) == 0 {
		str += "0"
	} else {
		str += string(i)
	}
	if len(f) != 0 {
		str += "." + string(f)
	}
	return str
}

func (n Number) ToFloat() float64 {
	v, err := strconv.ParseFloat(n.String(), 64)
	// Normally we shouldn't get an error here cause string value is valid.
	// But we should throw error just in case.
	if err != nil {
		panic(err)
	}
	return v
}
