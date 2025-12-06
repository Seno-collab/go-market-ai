package rbac

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	Repo UserRoleRepo
}

func (s Service) Check(ctx context.Context, userID uuid.UUID, perm string) (bool, error) {
	agg, err := s.Repo.GetUserRole(ctx, userID)
	if err != nil {
		return false, err
	}
	return agg.HasPermission(perm), nil
}
