package container

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/andvarfolomeev/docker-notifier/internal/docker"
)

type ClientOptions struct {
	LabelEnabled bool
}

type Client struct {
	SDK  DockerSDK
	Opts *ClientOptions
}

func NewClient(dockerSDK DockerSDK, opts *ClientOptions) (*Client, error) {
	return &Client{SDK: dockerSDK, Opts: opts}, nil
}

func (dc *Client) RunningContainers(ctx context.Context) ([]Container, error) {
	filterArgs := RunningContainerFilters(dc.Opts)
	dockerContainers, err := dc.SDK.ContainerList(ctx, docker.ContainerListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to list containers: %w", err)
	}

	containers := ConvertContainers(dockerContainers)
	return containers, nil
}

func (dc *Client) ContainerLogs(ctx context.Context, containerID, since string, tail int) ([]byte, error) {
	if ctx == nil {
		panic("context must not be nil")
	}

	if containerID == "" {
		return nil, errors.New("container ID cannot be empty")
	}

	logCtx, cancel := context.WithTimeout(ctx, PingTimeout)
	defer cancel()

	dockerOpts := ContainerLogsOptions(since, tail)
	logs, err := dc.SDK.ContainerLogs(logCtx, containerID, dockerOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for container %s: %w", containerID, err)
	}
	defer logs.Close()

	buf, err := io.ReadAll(logs)
	if err != nil {
		return nil, fmt.Errorf("failed to read logs for container %s: %w", containerID, err)
	}

	return buf, nil
}

func (dc *Client) Close() error {
	if dc.SDK != nil {
		dc.SDK.Close()
	}
	return nil
}
