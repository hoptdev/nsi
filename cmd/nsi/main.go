package main

import (
	"context"
	"fmt"
	"log/slog"
	"nsi/internal/app"
	"nsi/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envDev  = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.Load()

	log := setupLogger(cfg.Env)

	log.Info("start application", slog.String("mode", cfg.Env))

	application := app.New(log, cfg)

	go application.HttpServer.Run()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.HttpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
