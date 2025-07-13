package main

import (
	"context"
	"minivault/config"
	"minivault/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load() // Load .env file if present, ignore error if missing
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	server.Run(ctx, cfg)
}
