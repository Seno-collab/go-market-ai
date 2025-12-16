package utils

import (
	domainerr "go-ai/pkg/domain_err"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[\w\-.]+@([\w\-]+\.)+[\w\-]{2,4}$`)
	urlRegex   = regexp.MustCompile(`^https?://.+`)
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

// Url
type Url struct {
	value string
}

func NewUrl(v string) (Url, error) {
	if v == "" {
		return Url{value: ""}, nil
	}
	if !urlRegex.MatchString(v) {
		return Url{}, domainerr.ErrInvalidUrl
	}
	return Url{value: v}, nil
}

func (u Url) String() string {
	return u.value
}

type Money int64 // VND

func NewMoney(v int64) (Money, error) {
	if v < 0 {
		return 0, domainerr.ErrInvalidPrice
	}
	return Money(v), nil
}

func (m Money) Add(v Money) Money {
	return m + v
}
