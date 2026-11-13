package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"pwmanager/internal/models"
	"pwmanager/internal/repository"
)

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService(repo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func hashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}

func (s *APIKeyService) Create(ctx context.Context, userID int64, in models.APIKeyInput) (*models.APIKeyCreated, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, errors.New("key name is required")
	}
	if len(name) > 100 {
		return nil, errors.New("key name too long")
	}

	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	plain := "pwm_" + base64.RawURLEncoding.EncodeToString(buf)
	prefix := plain[:12]

	k, err := s.repo.Create(ctx, userID, name, hashAPIKey(plain), prefix)
	if err != nil {
		return nil, err
	}
	return &models.APIKeyCreated{APIKey: *k, Key: plain}, nil
}

func (s *APIKeyService) List(ctx context.Context, userID int64) ([]models.APIKey, error) {
	return s.repo.List(ctx, userID)
}

func (s *APIKeyService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}

func (s *APIKeyService) Validate(ctx context.Context, plain string) (int64, string, error) {
	if !strings.HasPrefix(plain, "pwm_") {
		return 0, "", repository.ErrNotFound
	}
	hash := hashAPIKey(plain)
	userID, username, keyID, err := s.repo.FindUserByHash(ctx, hash)
	if err != nil {
		return 0, "", err
	}
	go s.repo.Touch(context.Background(), keyID)
	return userID, username, nil
}
