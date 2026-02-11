package authapp

import (
	"context"

	"go-ai/internal/identity/domain/auth"
	domainerr "go-ai/pkg/domain_err"
	"go-ai/pkg/helpers"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

type RegisterUseCase struct {
	Repo auth.Repository
}

func NewRegisterUseCase(repo auth.Repository) *RegisterUseCase {
	return &RegisterUseCase{
		Repo: repo,
	}
}

func (s *RegisterUseCase) Execute(ctx context.Context, request RegisterRequest) (uuid.UUID, error) {

	email, err := helpers.NewEmail(request.Email)
	if err != nil {
		return uuid.Nil, err
	}

	rawPassword, err := auth.NewPassword(request.Password)
	if err != nil {
		return uuid.Nil, err
	}

	hashedPassword, err := helpers.HashPassword(rawPassword.String())
	if err != nil {
		return uuid.Nil, domainerr.ErrInternalServerError
	}
	pw, _ := auth.NewPasswordFromHash(hashedPassword)
	return s.Repo.CreateUser(ctx, &auth.Entity{
		FullName: request.FullName,
		Email:    email,
		Password: pw,
	})
}
