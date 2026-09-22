package users

import (
	"context"
	"errors"
	"strings"
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

	user := &User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}

	if err := s.repository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
