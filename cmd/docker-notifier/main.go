package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yourusername/docker-notifier/internal/config"
	"github.com/yourusername/docker-notifier/internal/logger"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		os.Exit(1)
	}

	log := logger.New(cfg.Debug)
	log.Info("🐳 Docker Notifier starting...")

	log.Info("Configuration:")
	log.Info("  Polling interval: %d seconds", cfg.Interval)
	log.Info("  Label filtering: %v", cfg.LabelEnable)
	log.Info("  Error patterns: %v", cfg.ErrorPatterns)
	log.Info("  Debug mode: %v", cfg.Debug)
	log.Info("  Cleanup mode: %v", cfg.Cleanup)

	// TODO: Initialize components
	// - Docker client
	// - Telegram client
	// - Watcher service

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal
	sig := <-sigChan
	log.Info("Received signal %v, shutting down...", sig)

	// TODO: Perform cleanup
}
