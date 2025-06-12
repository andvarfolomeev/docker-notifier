package container

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

type DockerSDK interface {
	ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	Ping(ctx context.Context) (types.Ping, error)
	Close() error
}
