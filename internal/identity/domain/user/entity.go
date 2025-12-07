package user

import (
	"go-ai/pkg/utils"

	"github.com/google/uuid"
)

type Entity struct {
	ID       uuid.UUID
	FullName string
	Email    utils.Email
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
