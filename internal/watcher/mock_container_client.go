package watcher

import (
	"context"
	"sync"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
)

// MockContainerClient is a mock implementation of ContainerClient for testing
type MockContainerClient struct {
	mu                  sync.Mutex
	containers          []container.Container
	logs                map[string][]byte
	containersErr       error
	logsErr             error
	closeCallCount      int
	containersCallCount int
	logsCallCount       int
}

// NewMockContainerClient creates a new MockContainerClient
func NewMockContainerClient() *MockContainerClient {
	return &MockContainerClient{
		logs: make(map[string][]byte),
	}
}

// SetContainers sets the containers to be returned by RunningContainers
func (m *MockContainerClient) SetContainers(containers []container.Container) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.containers = containers
}

// SetContainersError sets the error to be returned by RunningContainers
func (m *MockContainerClient) SetContainersError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.containersErr = err
}

// SetLogs sets the logs to be returned by ContainerLogs for a specific container
func (m *MockContainerClient) SetLogs(id string, logs []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs[id] = logs
}

// SetLogsError sets the error to be returned by ContainerLogs
func (m *MockContainerClient) SetLogsError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logsErr = err
}

// RunningContainers implements ContainerClient.RunningContainers
func (m *MockContainerClient) RunningContainers(ctx context.Context) ([]container.Container, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.containersCallCount++
	if m.containersErr != nil {
		return nil, m.containersErr
	}
	return m.containers, nil
}

// ContainerLogs implements ContainerClient.ContainerLogs
func (m *MockContainerClient) ContainerLogs(ctx context.Context, id, since string, tail int) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logsCallCount++
	if m.logsErr != nil {
		return nil, m.logsErr
	}
	logs, ok := m.logs[id]
	if !ok {
		return []byte{}, nil
	}
	return logs, nil
}

// Close implements ContainerClient.Close
func (m *MockContainerClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeCallCount++
	return nil
}

// ContainersCallCount returns the number of calls to RunningContainers
func (m *MockContainerClient) ContainersCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.containersCallCount
}

// LogsCallCount returns the number of calls to ContainerLogs
func (m *MockContainerClient) LogsCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.logsCallCount
}

// CloseCallCount returns the number of calls to Close
func (m *MockContainerClient) CloseCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closeCallCount
}
