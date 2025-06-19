package container

import (
	"fmt"

	"github.com/andvarfolomeev/docker-notifier/internal/docker"
)

func RunningContainerFilters(opts *ClientOptions) docker.Filters {
	filterArgs := docker.NewFilter()
	filterArgs.Add("status", "running")

	if opts.LabelEnabled {
		filterArgs.Add("label", fmt.Sprintf("%s=%s", LabelEnableKey, LabelEnableValue))
	}

	return *filterArgs
}

func ContainerLogsOptions(since string, tail int) docker.ContainerLogsOptions {
	var tailStr string
	if tail > 0 {
		tailStr = fmt.Sprintf("%d", tail)
	}

	opts := docker.ContainerLogsOptions{
		Stdout:    true,
		Stderr:    true,
		Since:     since,
		Timestamp: true,
		Follow:    false,
		Tail:      tailStr,
	}

	return opts
}
