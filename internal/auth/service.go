package auth

import (
	"context"
	"errors"
	"time"

	"eldercare/internal/users"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	userRepository *users.Repository
	jwtSecret      string
}

func NewService(
	userRepository *users.Repository,
	jwtSecret string,
) *Service {
	return &Service{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	)

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	signedToken, err := token.SignedString(
		[]byte(s.jwtSecret),
	)

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return signedToken, nil
}
