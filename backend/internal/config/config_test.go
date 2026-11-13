package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"pwmanager/internal/config"
)

const validYAML = `
server:
  host: "127.0.0.1"
  port: 9090
  jwt_secret: "secret"
  jwt_ttl_hours: 2
  encryption_key: "0123456789abcdef0123456789abcdef"
  read_timeout_seconds: 10
  write_timeout_seconds: 20
database:
  host: "localhost"
  port: 5432
  user: "u"
  password: "p"
  name: "d"
  sslmode: "disable"
  max_conns: 5
  connect_retries: 3
  connect_retry_delay_seconds: 1
docker:
  enabled: false
  compose_file: "./x.yaml"
cors:
  allowed_origins:
    - "http://localhost:3000"
`

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	cfg, err := config.Load(writeTmp(t, validYAML))
	if err != nil { t.Fatalf("Load: %v", err) }
	if cfg.Server.Port != 9090 || cfg.Server.JWTTTLHours != 2 {
		t.Errorf("cfg: %+v", cfg.Server)
	}
	if len(cfg.CORS.AllowedOrigins) != 1 {
		t.Errorf("cors: %v", cfg.CORS.AllowedOrigins)
	}
}

func TestLoad_Missing(t *testing.T) {
	if _, err := config.Load(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Error("expected error")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	if _, err := config.Load(writeTmp(t, "not: [valid:")); err == nil {
		t.Error("expected error")
	}
}

func TestLoad_ShortKey(t *testing.T) {
	if _, err := config.Load(writeTmp(t, "server:\n  encryption_key: \"short\"\n")); err == nil {
		t.Error("short key not caught")
	}
}

func TestDurationHelpers(t *testing.T) {
	cfg, _ := config.Load(writeTmp(t, validYAML))
	if cfg.ReadTimeout() != 10*time.Second { t.Error("ReadTimeout") }
	if cfg.WriteTimeout() != 20*time.Second { t.Error("WriteTimeout") }
	if cfg.JWTTTL() != 2*time.Hour { t.Error("JWTTTL") }
}
