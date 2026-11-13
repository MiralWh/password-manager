package docker_test

import (
	"context"
	"testing"

	"pwmanager/internal/docker"
)

func TestEnsureUp_MissingFile(t *testing.T) {
	err := docker.EnsureUp(context.Background(), "/nonexistent/compose.yaml")
	if err == nil {
		t.Error("expected error for missing compose file")
	}
}
