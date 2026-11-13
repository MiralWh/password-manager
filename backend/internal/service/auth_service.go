package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"pwmanager/internal/auth"
	"pwmanager/internal/models"
	"pwmanager/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	users *repository.UserRepository
	jwt   *auth.JWTManager
}

func NewAuthService(users *repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

func (s *AuthService) Login(ctx context.Context, in models.LoginInput) (*models.LoginResponse, error) {
	user, err := s.users.FindByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)) != nil {
		return nil, ErrInvalidCredentials
	}
	token, exp, err := s.jwt.Generate(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{Token: token, ExpiresAt: exp, Username: user.Username}, nil
}

func (s *AuthService) Register(ctx context.Context, in models.RegisterInput) (*models.LoginResponse, error) {
	if len(in.Username) < 3 || len(in.Password) < 6 {
		return nil, errors.New("username must be >=3 and password >=6 chars")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user, err := s.users.Create(ctx, in.Username, string(hash))
	if err != nil {
		return nil, ErrUserExists
	}
	token, exp, err := s.jwt.Generate(user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return &models.LoginResponse{Token: token, ExpiresAt: exp, Username: user.Username}, nil
}
