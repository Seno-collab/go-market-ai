package auth

import (
	"regexp"
	"strings"
)

type Password struct {
	value string
}

var (
	lowercase   = regexp.MustCompile(`[a-z]`)
	uppercase   = regexp.MustCompile(`[A-Z]`)
	number      = regexp.MustCompile(`[0-9]`)
	specialChar = regexp.MustCompile(`[@$!%*?&]`)
)

func NewPassword(v string) (Password, error) {
	if strings.TrimSpace(v) == "" {
		return Password{}, ErrPasswordTooShort
	}
	if len(v) < 6 {
		return Password{}, ErrPasswordTooShort
	}
	if !lowercase.MatchString(v) {
		return Password{}, ErrWeakPassword
	}
	if !uppercase.MatchString(v) {
		return Password{}, ErrWeakPassword
	}
	if !number.MatchString(v) {
		return Password{}, ErrWeakPassword
	}
	if !specialChar.MatchString(v) {
		return Password{}, ErrWeakPassword
	}
	return Password{value: v}, nil
}

func (p Password) String() string {
	return p.value
}

func NewPasswordFromHash(hash string) (Password, error) {
	return Password{value: hash}, nil
}
