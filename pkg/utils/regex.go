package utils

import (
	domainerr "go-ai/pkg/domain_err"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[\w\-.]+@([\w\-]+\.)+[\w\-]{2,4}$`)
)

// Email
type Email struct {
	value string
}

func NewEmail(v string) (Email, error) {
	if v == "" {
		return Email{value: ""}, nil
	}
	if !emailRegex.MatchString(v) {
		return Email{}, domainerr.ErrInvalidEmail
	}
	return Email{value: v}, nil
}

func (e Email) String() string {
	return e.value
}
