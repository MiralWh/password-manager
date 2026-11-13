package service_test

import (
	"context"
	"errors"
	"testing"

	"pwmanager/internal/crypto"
	"pwmanager/internal/models"
	"pwmanager/internal/repository"
	"pwmanager/internal/service"
	"pwmanager/internal/testutil"
)

const encKey = "0123456789abcdef0123456789abcdef"

func newPasswordSvc(t *testing.T) (*service.PasswordService, int64) {
	pool, _ := testutil.NewTestDB(t)
	repo := repository.NewPasswordRepository(pool)
	c, _ := crypto.New(encKey)
	return service.NewPasswordService(repo, c), testutil.FindUserID(t, pool, "admin")
}

func TestPasswordService_CreateDecryptsInResponse(t *testing.T) {
	svc, uid := newPasswordSvc(t)
	p, err := svc.Create(context.Background(), uid, models.PasswordInput{
		Title: "T", Username: "U", Password: "Plain",
	})
	if err != nil { t.Fatalf("Create: %v", err) }
	if p.Password != "Plain" { t.Errorf("plain: %q", p.Password) }
}

func TestPasswordService_Create_Validation(t *testing.T) {
	svc, uid := newPasswordSvc(t)
	for _, in := range []models.PasswordInput{
		{Username: "u", Password: "p"},
		{Title: "t", Password: "p"},
		{Title: "t", Username: "u"},
	} {
		if _, err := svc.Create(context.Background(), uid, in); err == nil {
			t.Errorf("expected validation error")
		}
	}
}

func TestPasswordService_CRUD(t *testing.T) {
	svc, uid := newPasswordSvc(t)
	ctx := context.Background()

	p, _ := svc.Create(ctx, uid, models.PasswordInput{Title: "A", Username: "B", Password: "p1"})

	upd, err := svc.Update(ctx, uid, p.ID, models.PasswordInput{Title: "A2", Username: "B2", Password: "p2"})
	if err != nil { t.Fatalf("Update: %v", err) }
	if upd.Title != "A2" || upd.Password != "p2" { t.Errorf("upd: %+v", upd) }

	if err := svc.Delete(ctx, uid, p.ID); err != nil { t.Fatalf("Delete: %v", err) }
	if _, err := svc.Get(ctx, uid, p.ID); !errors.Is(err, repository.ErrNotFound) {
		t.Error("not found after delete")
	}
}

func TestPasswordService_List_Search(t *testing.T) {
	svc, uid := newPasswordSvc(t)
	ctx := context.Background()
	svc.Create(ctx, uid, models.PasswordInput{Title: "Gmail", Username: "me", Password: "x"})
	svc.Create(ctx, uid, models.PasswordInput{Title: "Slack", Username: "work", Password: "y"})

	list, _ := svc.List(ctx, uid, "gmail")
	if len(list) != 1 || list[0].Title != "Gmail" { t.Errorf("search: %+v", list) }
	list, _ = svc.List(ctx, uid, "")
	if len(list) != 2 { t.Errorf("all: %d", len(list)) }
}
