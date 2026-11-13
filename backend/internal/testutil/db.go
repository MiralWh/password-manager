package testutil

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	"pwmanager/internal/config"
	"pwmanager/internal/database"
)

var (
	sharedPool      *pgxpool.Pool
	sharedCfg       *config.Config
	sharedContainer testcontainers.Container
	once            sync.Once
	initErr         error
)

func RunTestMain(m *testing.M) {
	code := m.Run()
	if sharedPool != nil {
		sharedPool.Close()
	}
	if sharedContainer != nil {
		_ = sharedContainer.Terminate(context.Background())
	}
	os.Exit(code)
}

func setup() {
	ctx := context.Background()
	pg, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("pwm_test"),
		postgres.WithUsername("pwm"),
		postgres.WithPassword("pwm_pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		initErr = fmt.Errorf("start container: %w", err)
		return
	}
	sharedContainer = pg

	host, err := pg.Host(ctx)
	if err != nil {
		initErr = err
		return
	}
	port, err := pg.MappedPort(ctx, "5432/tcp")
	if err != nil {
		initErr = err
		return
	}
	portNum := port.Int() // ← было port.Num() (ошибка)

	cfg := &config.Config{
		Server: config.ServerConfig{
			JWTSecret:           "test-secret-long-enough-for-hs256-32+",
			JWTTTLHours:         1,
			EncryptionKey:       "0123456789abcdef0123456789abcdef",
			ReadTimeoutSeconds:  15,
			WriteTimeoutSeconds: 15,
		},
		Database: config.DatabaseConfig{
			Host:                     host,
			Port:                     portNum,
			User:                     "pwm",
			Password:                 "pwm_pass",
			Name:                     "pwm_test",
			SSLMode:                  "disable",
			MaxConns:                 5,
			ConnectRetries:           20,
			ConnectRetryDelaySeconds: 1,
		},
		CORS: config.CORSConfig{AllowedOrigins: []string{"http://localhost:3000"}},
	}
	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		initErr = err
		return
	}
	if err := database.Migrate(ctx, pool); err != nil {
		initErr = err
		return
	}
	sharedPool = pool
	sharedCfg = cfg
}

func NewTestDB(t *testing.T) (*pgxpool.Pool, *config.Config) {
	t.Helper()
	once.Do(setup)
	if initErr != nil {
		t.Fatalf("testcontainers: %v (Docker Desktop запущен?)", initErr)
	}
	ctx := context.Background()
	if _, err := sharedPool.Exec(ctx,
		`TRUNCATE api_keys, passwords, groups, users RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	SeedUser(t, sharedPool, "admin", "admin123")
	return sharedPool, sharedCfg
}

func SeedUser(t *testing.T, pool *pgxpool.Pool, username, password string) int64 {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	var id int64
	err = pool.QueryRow(context.Background(),
		`INSERT INTO users (username, password_hash) VALUES ($1,$2) RETURNING id`,
		username, string(hash)).Scan(&id)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	return id
}

func FindUserID(t *testing.T, pool *pgxpool.Pool, username string) int64 {
	t.Helper()
	var id int64
	if err := pool.QueryRow(context.Background(),
		`SELECT id FROM users WHERE username=$1`, username).Scan(&id); err != nil {
		t.Fatalf("find user: %v", err)
	}
	return id
}
