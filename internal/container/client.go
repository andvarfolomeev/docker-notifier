package container

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ClientOptions struct {
	LabelEnabled bool
}

type Client struct {
	SDK  DockerSDK
	Opts *ClientOptions
}

func NewClient(opts *ClientOptions) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to init docker client: %w", err)
	}

	return &Client{SDK: cli, Opts: opts}, nil
}

func (dc *Client) RunningContainers(ctx context.Context) ([]Container, error) {
	filterArgs := RunningContainerFilters(dc.Opts)
	dockerContainers, err := dc.SDK.ContainerList(ctx, types.ContainerListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to list containers: %w", err)
	}

	containers := ConvertContainers(dockerContainers)
	return containers, nil
}

func (dc *Client) ContainerLogs(ctx context.Context, id, since string, tail int) ([]byte, error) {
	if ctx == nil {
		panic("context must not be nil")
	}

	if id == "" {
		return nil, errors.New("container ID cannot be empty")
	}

	logCtx, cancel := context.WithTimeout(ctx, PingTimeout)
	defer cancel()

	dockerOpts := ContainerLogsOptions(since, tail)
	logs, err := dc.SDK.ContainerLogs(logCtx, id, dockerOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for container %s: %w", id, err)
	}
	defer logs.Close()

	buf, err := io.ReadAll(logs)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs for container %s: %w", id, err)
	}

	return buf, nil
}

func (dc *Client) Close() error {
	if dc.SDK != nil {
		return dc.SDK.Close()
	}
	return nil
}
