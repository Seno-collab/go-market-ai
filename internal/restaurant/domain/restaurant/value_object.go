package restaurant

import (
	"regexp"
)

var (
	phoneRegex = regexp.MustCompile(`^[0-9]{9,15}$`)
)

// Phone
type Phone struct {
	value string
}

func NewPhone(v string) (Phone, error) {
	if v == "" {
		return Phone{value: ""}, nil
	}
	if !phoneRegex.MatchString(v) {
		return Phone{}, ErrInvalidPhone
	}
	return Phone{value: v}, nil
}

func (p Phone) String() string {
	return p.value
}
