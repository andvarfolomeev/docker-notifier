package container

import (
	"context"
	"io"

	"github.com/andvarfolomeev/docker-notifier/internal/docker"
)

type DockerSDK interface {
	ContainerList(context.Context, docker.ContainerListOptions) ([]docker.Container, error)
	ContainerLogs(context.Context, string, docker.ContainerLogsOptions) (io.ReadCloser, error)
	Ping(context.Context) (string, error)
	Close()

	// ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	// ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	// Ping(ctx context.Context) (types.Ping, error)
	// Close() error
}
