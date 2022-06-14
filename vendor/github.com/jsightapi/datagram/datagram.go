package datagram

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
)

const (
	maxSize = 508
)

type Datagram struct {
	index int
	list  []*element
}

func New() Datagram {
	return Datagram{
		list: make([]*element, 0, 10), //nolint:gomnd
	}
}

func (d *Datagram) Append(k, v string) {
	e := element{
		key:         k,
		value:       v,
		truncatable: false,
	}
	d.list = append(d.list, &e)
}

func (d *Datagram) AppendTruncatable(k, v string) {
	e := element{
		key:         k,
		value:       v,
		truncatable: true,
	}
	d.list = append(d.list, &e)
}

func (d *Datagram) Pack() []byte {
	for {
		s := d.pack()

		if len(s) <= maxSize {
			return s
		}

		e := d.truncatableElement()
		if e == nil {
			return s[:maxSize]
		}

		e.truncate(len(s) - maxSize)
	}
}

func (d Datagram) pack() []byte {
	b := new(bytes.Buffer)

	w := csv.NewWriter(b)
	_ = w.Write(d.keys())
	_ = w.Write(d.values())
	w.Flush()

	return b.Bytes()
}

// truncatableElement selects one element that can be truncated.
// Returns nil if not found.
func (d Datagram) truncatableElement() *element {
	for i := len(d.list) - 1; i >= 0; i-- {
		if d.list[i].truncatable {
			return d.list[i]
		}
	}
	return nil
}

func (d Datagram) keys() []string {
	ss := make([]string, 0, len(d.list))
	for i := 0; i < len(d.list); i++ {
		ss = append(ss, d.list[i].key)
	}
	return ss
}

func (d Datagram) values() []string {
	ss := make([]string, 0, len(d.list))
	for i := 0; i < len(d.list); i++ {
		ss = append(ss, d.list[i].value)
	}
	return ss
}

// Unpack the CSV. Replaces missing values with empty strings.
func Unpack(b []byte) (Datagram, error) {
	d := New()

	r := csv.NewReader(bytes.NewReader(b))

	keys, err := r.Read()
	if err == io.EOF {
		return d, errors.New("keys not found")
	}
	if err != nil {
		return d, err
	}

	values, err := r.Read()
	if err == io.EOF {
		return d, errors.New("values not found")
	}
	if err != nil && !errors.Is(err, csv.ErrFieldCount) {
		return d, err
	}

	for i := 0; i < len(keys); i++ {
		if i < len(values) {
			d.Append(keys[i], values[i])
		} else {
			d.Append(keys[i], "")
		}
	}

	return d, nil
}

func (d *Datagram) InitIterator() {
	d.index = 0
}

func (d Datagram) HasNext() bool {
	return d.index < len(d.list)
}

func (d *Datagram) GetNext() (k, v string) {
	if d.HasNext() {
		element := d.list[d.index]
		d.index++
		return element.key, element.value
	}
	return "", ""
}
