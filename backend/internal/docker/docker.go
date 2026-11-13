package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func EnsureUp(ctx context.Context, composeFile string) error {
	if _, err := os.Stat(composeFile); err != nil {
		return fmt.Errorf("compose file: %w", err)
	}
	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", composeFile, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose up: %w", err)
	}
	return nil
}
