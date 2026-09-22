package users

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(
	ctx context.Context,
	name string,
	email string,
	passwordHash string,
	role string,
) (*User, error) {

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	role = strings.TrimSpace(role)

	if name == "" {
		return nil, errors.New("name is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if passwordHash == "" {
		return nil, errors.New("password is required")
	}

	if role == "" {
		return nil, errors.New("role is required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(passwordHash), bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.New("Failed To hash password")
	}

	user := &User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
