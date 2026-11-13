package service

import (
	"context"
	"errors"

	"pwmanager/internal/crypto"
	"pwmanager/internal/models"
	"pwmanager/internal/repository"
)

type PasswordService struct {
	repo   *repository.PasswordRepository
	cipher *crypto.Cipher
}

func NewPasswordService(repo *repository.PasswordRepository, cipher *crypto.Cipher) *PasswordService {
	return &PasswordService{repo: repo, cipher: cipher}
}

func (s *PasswordService) decryptRow(row repository.PasswordRow) (models.Password, error) {
	dec, err := s.cipher.Decrypt(row.EncPassword)
	if err != nil {
		return models.Password{}, err
	}
	p := row.Password
	p.Password = string(dec)
	return p, nil
}

func (s *PasswordService) List(ctx context.Context, userID int64, search string) ([]models.Password, error) {
	rows, err := s.repo.ListByUser(ctx, userID, search)
	if err != nil {
		return nil, err
	}
	out := make([]models.Password, 0, len(rows))
	for _, r := range rows {
		p, err := s.decryptRow(r)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, nil
}

func (s *PasswordService) Get(ctx context.Context, userID, id int64) (*models.Password, error) {
	row, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	p, err := s.decryptRow(*row)
	return &p, err
}

func (s *PasswordService) Create(ctx context.Context, userID int64, in models.PasswordInput) (*models.Password, error) {
	if in.Title == "" || in.Username == "" || in.Password == "" {
		return nil, errors.New("title, username and password are required")
	}
	enc, err := s.cipher.Encrypt([]byte(in.Password))
	if err != nil {
		return nil, err
	}
	row, err := s.repo.Create(ctx, userID, in, enc)
	if err != nil {
		return nil, err
	}
	p, err := s.decryptRow(*row)
	return &p, err
}

func (s *PasswordService) Update(ctx context.Context, userID, id int64, in models.PasswordInput) (*models.Password, error) {
	if in.Title == "" || in.Username == "" || in.Password == "" {
		return nil, errors.New("title, username and password are required")
	}
	enc, err := s.cipher.Encrypt([]byte(in.Password))
	if err != nil {
		return nil, err
	}
	row, err := s.repo.Update(ctx, userID, id, in, enc)
	if err != nil {
		return nil, err
	}
	p, err := s.decryptRow(*row)
	return &p, err
}

func (s *PasswordService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}
