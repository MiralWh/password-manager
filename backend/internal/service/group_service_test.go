package service_test

import (
	"context"
	"errors"
	"testing"

	"pwmanager/internal/models"
	"pwmanager/internal/repository"
	"pwmanager/internal/service"
	"pwmanager/internal/testutil"
)

func newGroupSvc(t *testing.T) (*service.GroupService, int64) {
	pool, _ := testutil.NewTestDB(t)
	return service.NewGroupService(repository.NewGroupRepository(pool)),
		testutil.FindUserID(t, pool, "admin")
}

func TestGroupService_Validation(t *testing.T) {
	svc, uid := newGroupSvc(t)
	for _, in := range []models.GroupInput{
		{Name: ""},
		{Name: "   "},
		{Name: string(make([]byte, 200))},
	} {
		if _, err := svc.Create(context.Background(), uid, in); err == nil {
			t.Errorf("expected error for %+v", in)
		}
	}
}

func TestGroupService_CRUD(t *testing.T) {
	svc, uid := newGroupSvc(t)
	ctx := context.Background()

	g, err := svc.Create(ctx, uid, models.GroupInput{Name: "Работа"})
	if err != nil { t.Fatalf("Create: %v", err) }

	if _, err := svc.Rename(ctx, uid, g.ID, models.GroupInput{Name: "Job"}); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	if err := svc.Delete(ctx, uid, g.ID); err != nil { t.Fatalf("Delete: %v", err) }
	if err := svc.Delete(ctx, uid, g.ID); !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("second delete: %v", err)
	}
}
