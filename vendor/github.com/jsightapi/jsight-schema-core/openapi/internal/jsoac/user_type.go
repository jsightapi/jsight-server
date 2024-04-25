package jsoac

import (
	"fmt"
	"strings"
)

type UserType struct {
	name string
}

func (u UserType) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`{"$ref": "#/components/schemas/%s"}`, strings.TrimLeft(u.name, "@"))
	return []byte(s), nil
}
