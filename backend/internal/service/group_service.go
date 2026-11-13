package service

import (
	"context"
	"errors"
	"strings"

	"pwmanager/internal/models"
	"pwmanager/internal/repository"
)

type GroupService struct {
	repo *repository.GroupRepository
}

func NewGroupService(repo *repository.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) List(ctx context.Context, userID int64, search string) ([]models.Group, error) {
	return s.repo.List(ctx, userID, search)
}

func (s *GroupService) Create(ctx context.Context, userID int64, in models.GroupInput) (*models.Group, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, errors.New("group name is required")
	}
	if len(name) > 100 {
		return nil, errors.New("group name too long")
	}
	return s.repo.Create(ctx, userID, name)
}

func (s *GroupService) Rename(ctx context.Context, userID, id int64, in models.GroupInput) (*models.Group, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, errors.New("group name is required")
	}
	return s.repo.Rename(ctx, userID, id, name)
}

func (s *GroupService) Delete(ctx context.Context, userID, id int64) error {
	return s.repo.Delete(ctx, userID, id)
}
