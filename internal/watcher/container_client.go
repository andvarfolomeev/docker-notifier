package watcher

import (
	"context"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
)

type ContainerClient interface {
	RunningContainers(ctx context.Context) ([]container.Container, error)
	ContainerLogs(ctx context.Context, id, since string, tail int) ([]byte, error)
	Close() error
}
