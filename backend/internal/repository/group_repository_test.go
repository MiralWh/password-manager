package repository_test

import (
	"context"
	"errors"
	"testing"

	"pwmanager/internal/models"
	"pwmanager/internal/repository"
	"pwmanager/internal/testutil"
)

func TestGroupRepository_CRUD(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewGroupRepository(pool)
	ctx := context.Background()
	userID := testutil.FindUserID(t, pool, "admin")

	g, err := repo.Create(ctx, userID, "Work")
	if err != nil { t.Fatalf("Create: %v", err) }
	if g.ID == 0 || g.Name != "Work" { t.Fatalf("unexpected: %+v", g) }

	list, _ := repo.List(ctx, userID, "")
	if len(list) != 1 { t.Errorf("list: %d", len(list)) }

	if _, err := repo.Rename(ctx, userID, g.ID, "Job"); err != nil { t.Fatalf("Rename: %v", err) }
	list, _ = repo.List(ctx, userID, "Job")
	if len(list) != 1 { t.Error("search after rename failed") }

	if err := repo.Delete(ctx, userID, g.ID); err != nil { t.Fatalf("Delete: %v", err) }
	if _, err := repo.Rename(ctx, userID, g.ID, "x"); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestGroupRepository_UniquePerUser(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewGroupRepository(pool)
	ctx := context.Background()
	admin := testutil.FindUserID(t, pool, "admin")

	if _, err := repo.Create(ctx, admin, "Dup"); err != nil { t.Fatalf("first: %v", err) }
	if _, err := repo.Create(ctx, admin, "Dup"); err == nil {
		t.Error("UNIQUE not enforced")
	}

	other := testutil.SeedUser(t, pool, "other", "otherpw")
	if _, err := repo.Create(ctx, other, "Dup"); err != nil {
		t.Errorf("other user should create: %v", err)
	}
}

func TestGroupRepository_PasswordCount(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewGroupRepository(pool)
	pwRepo := repository.NewPasswordRepository(pool)
	ctx := context.Background()
	userID := testutil.FindUserID(t, pool, "admin")

	g, _ := repo.Create(ctx, userID, "Jobs")
	for i := 0; i < 3; i++ {
		pwRepo.Create(ctx, userID, models.PasswordInput{
			GroupID: &g.ID, Title: "T", Username: "U",
		}, []byte("e"))
	}

	list, _ := repo.List(ctx, userID, "")
	if len(list) != 1 || list[0].PasswordCount != 3 {
		t.Errorf("count: %+v", list)
	}
}

func TestGroupRepository_DeleteSetsNull(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	groupRepo := repository.NewGroupRepository(pool)
	pwRepo := repository.NewPasswordRepository(pool)
	ctx := context.Background()
	userID := testutil.FindUserID(t, pool, "admin")

	g, _ := groupRepo.Create(ctx, userID, "Temp")
	p, _ := pwRepo.Create(ctx, userID, models.PasswordInput{
		GroupID: &g.ID, Title: "T", Username: "U",
	}, []byte("e"))

	groupRepo.Delete(ctx, userID, g.ID)

	got, _ := pwRepo.Get(ctx, userID, p.ID)
	if got.GroupID != nil {
		t.Error("group_id should be NULL after group delete")
	}
}
