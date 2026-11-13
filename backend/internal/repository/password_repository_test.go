package repository_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"pwmanager/internal/models"
	"pwmanager/internal/repository"
	"pwmanager/internal/testutil"
)

func TestPasswordRepository_CRUD(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewPasswordRepository(pool)
	userID := testutil.FindUserID(t, pool, "admin")
	ctx := context.Background()

	enc := []byte("encrypted")
	row, err := repo.Create(ctx, userID, models.PasswordInput{
		Title: "GitHub", Username: "octocat", URL: "u", Notes: "n",
	}, enc)
	if err != nil { t.Fatalf("Create: %v", err) }
	if row.ID == 0 || !bytes.Equal(row.EncPassword, enc) {
		t.Fatalf("unexpected: %+v", row)
	}

	got, err := repo.Get(ctx, userID, row.ID)
	if err != nil { t.Fatalf("Get: %v", err) }

	newEnc := []byte("newenc")
	upd, err := repo.Update(ctx, userID, row.ID, models.PasswordInput{
		Title: "GH2", Username: "u2",
	}, newEnc)
	if err != nil { t.Fatalf("Update: %v", err) }
	if upd.Title != "GH2" || !bytes.Equal(upd.EncPassword, newEnc) {
		t.Errorf("update: %+v", upd)
	}

	list, _ := repo.ListByUser(ctx, userID, "")
	if len(list) != 1 { t.Errorf("list: %d", len(list)) }

	if err := repo.Delete(ctx, userID, row.ID); err != nil { t.Fatalf("Del: %v", err) }
	if _, err := repo.Get(ctx, userID, row.ID); !errors.Is(err, repository.ErrNotFound) {
		t.Error("not deleted")
	}
	_ = got
}

func TestPasswordRepository_Search(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewPasswordRepository(pool)
	ctx := context.Background()
	userID := testutil.FindUserID(t, pool, "admin")

	var groupID int64
	if err := pool.QueryRow(ctx,
		`INSERT INTO groups (user_id, name) VALUES ($1, 'Work') RETURNING id`, userID,
	).Scan(&groupID); err != nil { t.Fatal(err) }

	repo.Create(ctx, userID, models.PasswordInput{
		GroupID: &groupID, Title: "Gmail", Username: "me@gmail.com",
	}, []byte("x"))
	repo.Create(ctx, userID, models.PasswordInput{
		Title: "Slack", Username: "me@company.io",
	}, []byte("x"))

	// title
	list, _ := repo.ListByUser(ctx, userID, "gmail")
	if len(list) != 1 || list[0].Title != "Gmail" { t.Errorf("title search failed") }

	// username
	list, _ = repo.ListByUser(ctx, userID, "company")
	if len(list) != 1 || list[0].Title != "Slack" { t.Errorf("username search failed") }

	// group name
	list, _ = repo.ListByUser(ctx, userID, "Work")
	if len(list) != 1 || list[0].Title != "Gmail" { t.Errorf("group search failed") }

	// no match
	list, _ = repo.ListByUser(ctx, userID, "zzz")
	if len(list) != 0 { t.Errorf("expected empty") }
}

func TestPasswordRepository_Update_NotFound(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewPasswordRepository(pool)
	userID := testutil.FindUserID(t, pool, "admin")
	_, err := repo.Update(context.Background(), userID, 999999,
		models.PasswordInput{Title: "x"}, []byte("e"))
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestPasswordRepository_UserIsolation(t *testing.T) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewPasswordRepository(pool)
	ctx := context.Background()
	alice := testutil.SeedUser(t, pool, "alice", "alicepw")
	bob := testutil.SeedUser(t, pool, "bob", "bobpw")

	row, _ := repo.Create(ctx, alice, models.PasswordInput{Title: "T", Username: "U"}, []byte("e"))
	if _, err := repo.Get(ctx, bob, row.ID); !errors.Is(err, repository.ErrNotFound) {
		t.Error("user isolation broken")
	}
}
