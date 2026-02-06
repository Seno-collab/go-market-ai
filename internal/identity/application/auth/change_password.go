package authapp

import (
	"context"
	"go-ai/internal/identity/domain/auth"
	"go-ai/pkg/helpers"
	domainerr "go-ai/pkg/domain_err"

	"github.com/google/uuid"
)

type ChangePasswordRequest struct {
	NewPassword     string `json:"new_password"`
	OldPassword     string `json:"old_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ChangePasswordUseCase struct {
	Repo auth.Repository
}

func NewChangePasswordUseCase(repo auth.Repository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		Repo: repo,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, req ChangePasswordRequest, userID uuid.UUID) error {
	passwordHashCurrent, err := uc.Repo.GetPasswordByID(ctx, userID)
	if err != nil {
		return domainerr.ErrInternalServerError
	}
	if !helpers.CheckPasswordHash(req.OldPassword, passwordHashCurrent) {
		return auth.ErrOldPasswordIncorrect
	}
	if req.ConfirmPassword != req.NewPassword {
		return auth.ErrConfirmPassword
	}
	newPassword, err := auth.NewPassword(req.NewPassword)
	if err != nil {
		return err
	}
	hashPassword, err := helpers.HashPassword(newPassword.String())
	if err != nil {
		return auth.ErrHashPasswordFailed
	}
	return uc.Repo.ChangePassword(ctx, hashPassword, userID)
}
