package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"pwmanager/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.Name, cfg.Database.SSLMode)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	if cfg.Database.MaxConns > 0 {
		poolCfg.MaxConns = int32(cfg.Database.MaxConns)
	}

	var pool *pgxpool.Pool
	retries := cfg.Database.ConnectRetries
	delay := time.Duration(cfg.Database.ConnectRetryDelaySeconds) * time.Second

	for i := 0; i < retries; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, poolCfg)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			err = pool.Ping(pingCtx)
			cancel()
			if err == nil {
				return pool, nil
			}
			pool.Close()
		}
		fmt.Printf("db not ready (attempt %d/%d): %v\n", i+1, retries, err)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil, fmt.Errorf("could not connect to db after %d attempts: %w", retries, err)
}
