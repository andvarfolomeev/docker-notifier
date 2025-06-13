package container

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func RunningContainerFilters(opts *ClientOptions) filters.Args {
	filterArgs := filters.NewArgs()
	filterArgs.Add("status", "running")

	if opts.LabelEnabled {
		filterArgs.Add("label", fmt.Sprintf("%s=%s", LabelEnableKey, LabelEnableValue))
	}

	return filterArgs
}

func ContainerLogsOptions(since string, tail int) types.ContainerLogsOptions {
	var tailStr string
	if tail > 0 {
		tailStr = fmt.Sprintf("%d", tail)
	}

	opts := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Since:      since,
		Timestamps: true,
		Follow:     false,
		Tail:       tailStr,
	}

	return opts
}
