package authservice

import (
	"go-ai/internal/domain/auth"

	"github.com/google/uuid"
)

type Service struct {
	repo auth.Repository
}

func NewAuthService(repo auth.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(u *auth.Auth) (uuid.UUID, error) {
	if err := u.Validate(); err != nil {
		return uuid.Nil, err
	}
	// unique email
	if _, err := s.repo.GetByEmail(u.Email); err != nil {
		return uuid.Nil, auth.ErrConflict
	}
	return s.repo.CreateUser(u)
}

// }

// func (s *Service) Login(u Auth) (string, string) {

// }
