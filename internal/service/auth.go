package service

import (
	"context"
	"errors"

	"ChipHouse/internal/domain"
	"ChipHouse/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	users *repository.UserRepo
}

func NewAuthService(users *repository.UserRepo) *AuthService {
	return &AuthService{users: users}
}

func (a *AuthService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	u, err := a.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return u, nil
}

// Authorize checks whether a user's role is allowed for the given action.
func (a *AuthService) Authorize(u *domain.User, allowed ...domain.Role) bool {
	if u == nil {
		return false
	}
	for _, r := range allowed {
		if u.Role == r {
			return true
		}
	}
	return false
}
