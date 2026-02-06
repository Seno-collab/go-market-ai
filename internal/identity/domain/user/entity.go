package user

import (
	"go-ai/pkg/helpers"

	"github.com/google/uuid"
)

type Entity struct {
	ID       uuid.UUID
	FullName string
	Email    helpers.Email
	Role     string
	IsActive bool
}

// func (u *User) Validate() error {
// 	if !strings.Contains(u.Email, "@") {
// 		return ErrInvalidEmail
// 	}
// 	if strings.TrimSpace(u.FullName) == "" {
// 		return ErrInvalidName
// 	}
// 	return nil
// }
