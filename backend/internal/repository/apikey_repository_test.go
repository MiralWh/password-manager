package repository_test

import (
	"context"
	"testing"

	"pwmanager/internal/repository"
	"pwmanager/internal/testutil"
)

func TestAPIKeyRepository_CRUD(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewAPIKeyRepository(pool)
	ctx := context.Background()
	userID := testutil.FindUserID(t, pool, "admin")

	k, err := repo.Create(ctx, userID, "desktop", "hash1", "pwm_abcdefgh")
	if err != nil { t.Fatalf("Create: %v", err) }
	if k.Name != "desktop" || k.KeyPrefix != "pwm_abcdefgh" {
		t.Fatalf("unexpected: %+v", k)
	}

	list, _ := repo.List(ctx, userID)
	if len(list) != 1 { t.Errorf("list: %d", len(list)) }

	if err := repo.Delete(ctx, userID, k.ID); err != nil { t.Fatalf("Delete: %v", err) }
	list, _ = repo.List(ctx, userID)
	if len(list) != 0 { t.Error("not deleted") }
}

func TestAPIKeyRepository_FindAndTouch(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewAPIKeyRepository(pool)
	ctx := context.Background()
	userID := testutil.FindUserID(t, pool, "admin")

	repo.Create(ctx, userID, "n", "uniquehash", "pwm_xyz")
	uid, uname, keyID, err := repo.FindUserByHash(ctx, "uniquehash")
	if err != nil { t.Fatalf("Find: %v", err) }
	if uid != userID || uname != "admin" || keyID == 0 {
		t.Errorf("unexpected: %d %s %d", uid, uname, keyID)
	}

	repo.Touch(ctx, keyID)
	list, _ := repo.List(ctx, userID)
	if list[0].LastUsedAt == nil { t.Error("last_used_at not set") }
}

func TestAPIKeyRepository_DuplicateHash(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewAPIKeyRepository(pool)
	ctx := context.Background()
	userID := testutil.FindUserID(t, pool, "admin")

	repo.Create(ctx, userID, "a", "dup", "pwm_1")
	if _, err := repo.Create(ctx, userID, "b", "dup", "pwm_2"); err == nil {
		t.Error("dup hash not rejected")
	}
}
