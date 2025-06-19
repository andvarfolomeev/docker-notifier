package container

import (
	"strings"

	"github.com/andvarfolomeev/docker-notifier/internal/docker"
)

type Container struct {
	ID   string
	Name string
}

func ContainerName(container docker.Container) string {
	if len(container.Names) > 0 {
		name := container.Names[0]
		return strings.TrimPrefix(name, "/")
	}
	if len(container.ID) >= ShortIDLen {
		return container.ID[:ShortIDLen]
	}
	return container.ID
}

func ConvertContainers(dockerContainers []docker.Container) []Container {
	containers := make([]Container, 0, len(dockerContainers))
	for _, dockerContainer := range dockerContainers {
		containers = append(containers, Container{
			ID:   dockerContainer.ID,
			Name: ContainerName(dockerContainer),
		})
	}
	return containers
}
