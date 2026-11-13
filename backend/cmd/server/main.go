package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"pwmanager/internal/config"
	"pwmanager/internal/database"
	"pwmanager/internal/docker"
	"pwmanager/internal/server"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if cfg.Docker.Enabled {
		if err := docker.EnsureUp(ctx, cfg.Docker.ComposeFile); err != nil {
			log.Printf("warning: docker compose up failed: %v", err)
		}
	}
	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer db.Close()
	if err := database.Migrate(ctx, db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	if err := database.Seed(ctx, db); err != nil {
		log.Fatalf("seed: %v", err)
	}
	srv, err := server.New(cfg, db)
	if err != nil {
		log.Fatalf("server: %v", err)
	}
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("server: %v", err)
	}
}
