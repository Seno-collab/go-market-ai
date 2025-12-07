package restaurant

import (
	"regexp"
)

var (
	phoneRegex = regexp.MustCompile(`^[0-9]{9,15}$`)
	urlRegex   = regexp.MustCompile(`^https?://.+`)
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

// Url
type Url struct {
	value string
}

func NewUrl(v string) (Url, error) {
	if v == "" {
		return Url{value: ""}, nil
	}
	if !urlRegex.MatchString(v) {
		return Url{}, ErrInvalidUrl
	}
	return Url{value: v}, nil
}

func (u Url) String() string {
	return u.value
}
