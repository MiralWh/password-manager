package repository_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
	"pwmanager/internal/repository"
	"pwmanager/internal/testutil"
)

func TestUserRepository_CreateFind(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewUserRepository(pool)
	ctx := context.Background()

	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.MinCost)
	u, err := repo.Create(ctx, "bob", string(hash))
	if err != nil { t.Fatalf("Create: %v", err) }
	if u.ID == 0 || u.Username != "bob" { t.Fatalf("unexpected: %+v", u) }

	got, err := repo.FindByUsername(ctx, "bob")
	if err != nil { t.Fatalf("Find: %v", err) }
	if got.ID != u.ID || got.PasswordHash != string(hash) {
		t.Errorf("mismatch")
	}
}

func TestUserRepository_NotFound(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewUserRepository(pool)
	if _, err := repo.FindByUsername(context.Background(), "ghost"); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestUserRepository_Duplicate(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewUserRepository(pool)
	ctx := context.Background()
	repo.Create(ctx, "dup", "h")
	if _, err := repo.Create(ctx, "dup", "h"); err == nil {
		t.Error("duplicate not rejected")
	}
}
