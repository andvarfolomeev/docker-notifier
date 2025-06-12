package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/andvarfolomeev/docker-notifier/internal/config"
	"github.com/andvarfolomeev/docker-notifier/internal/container"
	"github.com/andvarfolomeev/docker-notifier/internal/telegram"
	"github.com/andvarfolomeev/docker-notifier/internal/watcher"
)

const (
	socketPermissionDeniedLog = "Permission denied accessing Docker socket. This is usually due to the container not having the correct permissions." +
		"Try running with: docker run -v /var/run/docker.sock:/var/run/docker.sock --group-add=$(stat -c '%g' /var/run/docker.sock) ..." +
		"Or use the provided scripts/start.sh script which handles permissions automatically."
	containerPermissionDeniedLog = "Permission denied when listing containers." +
		"Please check the Docker socket permissions and try again." +
		"If running in a container, make sure it has access to the Docker socket."
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		os.Exit(1)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	log := slog.New(handler)
	log.Info("🐳 Docker Notifier starting...")

	telegramClient := telegram.New(cfg.TelegramToken, cfg.TelegramChatID, log)

	containerClient, err := container.NewClient(&container.ClientOptions{
		LabelEnabled: cfg.LabelEnable,
	})

	if err != nil {
		if errors.Is(err, fs.ErrPermission) || strings.Contains(err.Error(), "permission denied") {
			log.Error(socketPermissionDeniedLog)
		} else {
			log.Error("Failed to initialize Docker client", "err", err)
		}
		os.Exit(1)
	}
	defer containerClient.Close()

	w, err := watcher.New(
		containerClient,
		&watcher.WatcherOptions{
			Interval:      time.Second * time.Duration(cfg.Interval),
			ErrorPatterns: cfg.ErrorPatterns,
		},
	)

	if err != nil {
		fmt.Println("Failed to initialize watcher", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// check permissions
	_, err = containerClient.RunningContainers(ctx)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) || strings.Contains(err.Error(), "permission denied") {
			log.Error(containerPermissionDeniedLog)
		} else {
			log.Error("Failed to list containers", "err", err)
		}
		cancel()
	}

	w.Start(ctx)

	go func() {
		for {
			select {
			case match, ok := <-w.C:
				if !ok {
					return
				}
				message := prepareMessage(match)
				err := telegramClient.SendMessage(ctx, message)
				if err != nil {
					log.Error("Failed to send message", "err", err)
				}
			case <-ctx.Done():
				log.Error("closed", "err", ctx.Err())
				return
			}
		}
	}()

	log.Info("Watcher started, polling logs", "interval", cfg.Interval)
	if cfg.LabelEnable {
		log.Info("Only containers with label %s=%s will be monitored", "com.andvarfolomeev.dockernotify.enable", "true")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("Received signal, shutting down...", "sig", sig)

	cancel()
	w.Cleanup()
}

func prepareMessage(match *watcher.MatchedLog) string {
	errorLine := match.Line.Content
	if len(errorLine) > 100 {
		errorLine = errorLine[:100]
	}

	messageLines := []string{
		fmt.Sprintf("🚨 Error detected!"),
		fmt.Sprintf("Container ID = %s; Container name = %s", match.Container.ID, match.Container.Name),
		fmt.Sprintf("Line: \"%s\"", errorLine),
	}
	message := strings.Join(messageLines, "\n")

	return message
}
