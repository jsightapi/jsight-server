package jsoac

type UserTypeArray struct {
	UserType
}

func (u UserTypeArray) MarshalJSON() ([]byte, error) {
	b, err := u.UserType.MarshalJSON()
	if err != nil {
		return b, err
	}

	b = append([]byte{'['}, b...)
	b = append(b, ']')

	return b, err
}
