package container

import (
	"strings"

	"github.com/docker/docker/api/types"
)

type Container struct {
	ID   string
	Name string
}

func containerName(container types.Container) string {
	if len(container.Names) > 0 {
		name := container.Names[0]
		return strings.TrimPrefix(name, "/")
	}
	if len(container.ID) >= shortIDLen {
		return container.ID[:shortIDLen]
	}
	return container.ID
}

func convertContainers(dockerContainers []types.Container) []Container {
	containers := make([]Container, 0, len(dockerContainers))
	for _, dockerContainer := range dockerContainers {
		containers = append(containers, Container{
			ID:   dockerContainer.ID,
			Name: containerName(dockerContainer),
		})
	}
	return containers
}
