package auth

import (
	"go-ai/pkg/utils"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var emailRegex = regexp.MustCompile(`^[\w\-.]+@([\w\-]+\.)+[\w\-]{2,4}$`)

type Entity struct {
	ID       uuid.UUID
	FullName string
	Email    utils.Email
	Password Password
	Role     string
	IsActive bool
}

func NewAuth(fullName, email, password, role string) (*Entity, error) {

	if err := validateFullName(fullName); err != nil {
		return nil, err
	}

	if err := validateEmail(email); err != nil {
		return nil, err
	}

	pw, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	if err := validateRole(role); err != nil {
		return nil, err
	}

	em, err := utils.NewEmail(email)
	if err != nil {
		return nil, ErrInvalidEmail
	}

	return &Entity{
		ID:       uuid.New(),
		FullName: fullName,
		Email:    em,
		Password: pw,
		Role:     role,
		IsActive: true,
	}, nil
}
func validateFullName(v string) error {
	if strings.TrimSpace(v) == "" {
		return ErrFullNameRequired
	}
	return nil
}

func validateEmail(v string) error {
	if !emailRegex.MatchString(v) {
		return ErrInvalidEmail
	}
	return nil
}

func validateRole(v string) error {
	if strings.TrimSpace(v) == "" {
		return ErrRoleRequired
	}
	return nil
}

func (e *Entity) Validate() error {
	if err := validateFullName(e.FullName); err != nil {
		return err
	}

	if err := validateEmail(e.Email.String()); err != nil {
		return err
	}

	if len(e.Password.String()) < 6 {
		return ErrPasswordTooShort
	}

	if err := validateRole(e.Role); err != nil {
		return err
	}

	return nil
}

func (e *Entity) UpdateEmail(v string) error {
	if err := validateEmail(v); err != nil {
		return err
	}
	em, err := utils.NewEmail(v)
	if err != nil {
		return ErrInvalidEmail
	}
	e.Email = em
	return nil
}

func (e *Entity) UpdatePassword(v string) error {
	pw, err := NewPassword(v)
	if err != nil {
		return err
	}
	e.Password = pw
	return nil
}

func (e *Entity) UpdateFullName(v string) error {
	if err := validateFullName(v); err != nil {
		return err
	}
	e.FullName = v
	return nil
}

func (e *Entity) UpdateRole(v string) error {
	if err := validateRole(v); err != nil {
		return err
	}
	e.Role = v
	return nil
}
