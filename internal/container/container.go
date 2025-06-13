package container

import (
	"strings"

	"github.com/docker/docker/api/types"
)

type Container struct {
	ID   string
	Name string
}

func ContainerName(container types.Container) string {
	if len(container.Names) > 0 {
		name := container.Names[0]
		return strings.TrimPrefix(name, "/")
	}
	if len(container.ID) >= ShortIDLen {
		return container.ID[:ShortIDLen]
	}
	return container.ID
}

func ConvertContainers(dockerContainers []types.Container) []Container {
	containers := make([]Container, 0, len(dockerContainers))
	for _, dockerContainer := range dockerContainers {
		containers = append(containers, Container{
			ID:   dockerContainer.ID,
			Name: ContainerName(dockerContainer),
		})
	}
	return containers
}
