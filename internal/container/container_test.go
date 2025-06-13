package container_test

import (
	"testing"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
	"github.com/docker/docker/api/types"
)

func TestContainerName(t *testing.T) {
	testCases := []struct {
		name         string
		container    types.Container
		expectedName string
	}{
		{
			name: "container with name",
			container: types.Container{
				ID:    "abcdef1234567890",
				Names: []string{"/test-container"},
			},
			expectedName: "test-container",
		},
		{
			name: "container with multiple names",
			container: types.Container{
				ID:    "abcdef1234567890",
				Names: []string{"/first-name", "/second-name"},
			},
			expectedName: "first-name",
		},
		{
			name: "container without name",
			container: types.Container{
				ID:    "abcdef1234567890",
				Names: []string{},
			},
			expectedName: "abcdef123456",
		},
		{
			name: "container with short ID",
			container: types.Container{
				ID:    "abcdef",
				Names: []string{},
			},
			expectedName: "abcdef",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			name := container.ContainerName(tc.container)
			if name != tc.expectedName {
				t.Errorf("expected name '%s', got '%s'", tc.expectedName, name)
			}
		})
	}
}

func TestConvertContainers(t *testing.T) {
	testCases := []struct {
		name               string
		dockerContainers   []types.Container
		expectedContainers []container.Container
	}{
		{
			name:               "empty list",
			dockerContainers:   []types.Container{},
			expectedContainers: []container.Container{},
		},
		{
			name: "single container",
			dockerContainers: []types.Container{
				{ID: "container1", Names: []string{"/test-container-1"}},
			},
			expectedContainers: []container.Container{
				{ID: "container1", Name: "test-container-1"},
			},
		},
		{
			name: "multiple containers",
			dockerContainers: []types.Container{
				{ID: "container1", Names: []string{"/test-container-1"}},
				{ID: "container2", Names: []string{"/test-container-2"}},
				{ID: "container3", Names: []string{"/test-container-3"}},
			},
			expectedContainers: []container.Container{
				{ID: "container1", Name: "test-container-1"},
				{ID: "container2", Name: "test-container-2"},
				{ID: "container3", Name: "test-container-3"},
			},
		},
		{
			name: "containers with different name formats",
			dockerContainers: []types.Container{
				{ID: "container1", Names: []string{"/test-container-1"}},
				{ID: "container2", Names: []string{}},
				{ID: "container3456789012", Names: []string{"/test-container-3"}},
			},
			expectedContainers: []container.Container{
				{ID: "container1", Name: "test-container-1"},
				{ID: "container2", Name: "container2"},
				{ID: "container3456789012", Name: "test-container-3"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			containers := container.ConvertContainers(tc.dockerContainers)

			if len(containers) != len(tc.expectedContainers) {
				t.Errorf("expected %d containers, got %d", len(tc.expectedContainers), len(containers))
			}

			for i, expected := range tc.expectedContainers {
				if i >= len(containers) {
					t.Errorf("missing expected container %d", i)
					continue
				}
				actual := containers[i]
				if actual.ID != expected.ID || actual.Name != expected.Name {
					t.Errorf("container %d mismatch: expected %+v, got %+v", i, expected, actual)
				}
			}
		})
	}
}
