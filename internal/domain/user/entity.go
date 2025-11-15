package user

import "strings"

type User struct {
	ID    int64
	Email string
	Name  string
}

func (u *User) Validate() error {
	if !strings.Contains(u.Email, "@") {
		return ErrInvalidEmail
	}
	if strings.TrimSpace(u.Name) == "" {
		return ErrInvalidName
	}
	return nil
}
