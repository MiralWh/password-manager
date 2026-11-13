package service_test

import (
	"context"
	"strings"
	"testing"

	"pwmanager/internal/models"
	"pwmanager/internal/repository"
	"pwmanager/internal/service"
	"pwmanager/internal/testutil"
)

func newAPIKeySvc(t *testing.T) (*service.APIKeyService, int64) {
	pool, _ := testutil.NewTestDB(t)
	return service.NewAPIKeyService(repository.NewAPIKeyRepository(pool)),
		testutil.FindUserID(t, pool, "admin")
}

func TestAPIKeyService_Create(t *testing.T) {
	svc, uid := newAPIKeySvc(t)
	k, err := svc.Create(context.Background(), uid, models.APIKeyInput{Name: "desktop"})
	if err != nil { t.Fatalf("Create: %v", err) }
	if !strings.HasPrefix(k.Key, "pwm_") { t.Errorf("bad prefix: %q", k.Key) }
	if len(k.Key) < 20 { t.Error("key too short") }
	if k.KeyPrefix != k.Key[:12] { t.Errorf("prefix mismatch") }

	list, _ := svc.List(context.Background(), uid)
	if list[0].KeyPrefix == k.Key {
		t.Error("plaintext leaked")
	}
}

func TestAPIKeyService_Validate(t *testing.T) {
	svc, uid := newAPIKeySvc(t)
	ctx := context.Background()
	k, _ := svc.Create(ctx, uid, models.APIKeyInput{Name: "k"})

	gotUID, gotName, err := svc.Validate(ctx, k.Key)
	if err != nil { t.Fatalf("Validate: %v", err) }
	if gotUID != uid || gotName != "admin" { t.Errorf("unexpected") }

	if _, _, err := svc.Validate(ctx, "wrong-prefix"); err == nil {
		t.Error("wrong prefix should fail")
	}
	if _, _, err := svc.Validate(ctx, "pwm_notexistingkey1234567890"); err == nil {
		t.Error("unknown key should fail")
	}
}

func TestAPIKeyService_Delete(t *testing.T) {
	svc, uid := newAPIKeySvc(t)
	ctx := context.Background()
	k, _ := svc.Create(ctx, uid, models.APIKeyInput{Name: "k"})
	if err := svc.Delete(ctx, uid, k.ID); err != nil { t.Fatalf("Delete: %v", err) }
	list, _ := svc.List(ctx, uid)
	if len(list) != 0 { t.Error("not deleted") }
}
