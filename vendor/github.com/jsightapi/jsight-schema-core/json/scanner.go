package json

import (
	"github.com/jsightapi/jsight-schema-core/bytes"
	"github.com/jsightapi/jsight-schema-core/errs"
)

type scanner struct {
	stateFn func(byte) bool

	index    int
	intLen   int
	fraLen   int
	expBegin int

	finished bool
	negative bool
}

func newScanner() *scanner {
	s := &scanner{}
	s.stateFn = s.stateOnSearchStart
	return s
}

func (s *scanner) Scan(value bytes.Bytes) (*Number, error) {
	for i, c := range value.Data() {
		s.index = i
		s.finished = true
		if !s.stateFn(c) {
			return nil, errs.ErrIncorrectNumberValue.F(value.String())
		}
	}

	if !s.finished {
		return nil, errs.ErrIncorrectNumberValue.F(value.String())
	}

	if err := s.setExp(value); err != nil {
		return nil, err
	}

	n := Number{
		neg: s.negative,
		nat: s.getNatural(value),
		exp: s.fraLen,
	}

	err := n.trimLeadingZerosInTheIntegerPart()
	if err != nil {
		return nil, err
	}

	err = n.trimTrailingZerosInTheFractionalPart()
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (s *scanner) setExp(value bytes.Bytes) error {
	if s.expBegin == 0 {
		return nil
	}

	exp, err := value.SubLow(s.expBegin).ParseInt()
	if err != nil {
		return err
	}
	// example with negative exp: 12.34E-1 = 1.234; exp = -1; intLen = 2 + (-1) = 1
	// example with positive exp: 12.34E+1 = 123.4; exp =  1; intLen = 2 + 1    = 3
	s.intLen += exp
	s.fraLen -= exp
	return nil
}

func (s *scanner) getNatural(value bytes.Bytes) bytes.Bytes {
	var natural bytes.Bytes

	switch {
	case s.intLen < 0: // example 1.2E-2 = .012
		natural = bytes.MakeBytes(s.fraLen)
		natural = appendZeros(natural, -s.intLen)
		natural = appendDigits(value, natural)

	case s.fraLen < 0: // example 1.2E+2 = 120
		natural = bytes.MakeBytes(s.intLen)
		natural = appendDigits(value, natural)
		natural = appendZeros(natural, -s.fraLen)
		s.fraLen = 0

	default: // example 12.3E-1 = 1.23
		natural = bytes.MakeBytes(s.intLen + s.fraLen)
		natural = appendDigits(value, natural)
	}

	return natural
}

func (s *scanner) stateOnSearchStart(c byte) bool {
	switch c {
	case '-':
		s.negative = true
		s.finished = false
		s.stateFn = s.stateMinusFound

	case '0':
		s.intLen++
		s.stateFn = s.stateFirstZeroFound

	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.intLen++
		s.stateFn = s.stateIntegerNumberFound
	default:
		return false
	}
	return true
}

func (s *scanner) stateMinusFound(c byte) bool {
	switch c {
	case '0':
		s.intLen++
		s.stateFn = s.stateFirstZeroFound
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.intLen++
		s.stateFn = s.stateIntegerNumberFound
	default:
		return false
	}
	return true
}

func (s *scanner) stateFirstZeroFound(c byte) bool {
	if c == '.' {
		s.stateFn = s.statePointFound
		return true
	}
	return false
}

func (s *scanner) stateIntegerNumberFound(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.intLen++

	case '.':
		s.stateFn = s.statePointFound

	case 'e', 'E':
		s.stateFn = s.stateExpFound
	default:
		return false
	}
	return true
}

func (s *scanner) statePointFound(c byte) bool {
	if '0' <= c && c <= '9' {
		s.fraLen++
		s.stateFn = s.stateFractionalNumberFound
		return true
	}
	return false
}

func (s *scanner) stateFractionalNumberFound(c byte) bool {
	switch c {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.fraLen++
	case 'e', 'E':
		s.stateFn = s.stateExpFound
	default:
		return false
	}
	return true
}

func (s *scanner) stateExpFound(c byte) bool {
	switch c {
	case '+':
		s.stateFn = s.stateExpSignFound

	case '-':
		if s.expBegin == 0 {
			s.expBegin = s.index
		}
		s.stateFn = s.stateExpSignFound

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		if s.expBegin == 0 {
			s.expBegin = s.index
		}
	default:
		return false
	}
	return true
}

func (s *scanner) stateExpSignFound(c byte) bool {
	if '0' <= c && c <= '9' {
		if s.expBegin == 0 {
			s.expBegin = s.index
		}
		s.stateFn = s.stateExpNumberFound
		return true
	}
	return false
}

func (*scanner) stateExpNumberFound(c byte) bool {
	return '0' <= c && c <= '9'
}

func appendZeros(to bytes.Bytes, n int) bytes.Bytes {
	for ; n > 0; n-- {
		to.Append('0')
	}
	return to
}

func appendDigits(from, to bytes.Bytes) bytes.Bytes {
loop:
	for _, c := range from.Data() {
		switch c {
		case '-', '.':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			to.Append(c)
		default:
			break loop
		}
	}
	return to
}
