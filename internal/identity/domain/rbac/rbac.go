package rbac

import "github.com/google/uuid"

type Permission string

type Role struct {
	ID         int
	Permission string
}

type UserRole struct {
	UserID uuid.UUID
	Role   string
}

func (ur UserRole) HasPermission(p string) bool {
	if ur.Role == p {
		return true
	}
	return false
}
