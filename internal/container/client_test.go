package container_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
	"github.com/docker/docker/api/types"
)

type mockDockerSDK struct {
	containerListFunc  func(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error)
	containerLogsFunc  func(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error)
	pingFunc           func(ctx context.Context) (types.Ping, error)
	closeFunc          func() error
	containerListCalls int
	containerLogsCalls int
	pingCalls          int
	closeCalls         int
}

func (m *mockDockerSDK) ContainerList(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
	m.containerListCalls++
	return m.containerListFunc(ctx, options)
}

func (m *mockDockerSDK) ContainerLogs(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
	m.containerLogsCalls++
	return m.containerLogsFunc(ctx, container, options)
}

func (m *mockDockerSDK) Ping(ctx context.Context) (types.Ping, error) {
	m.pingCalls++
	return m.pingFunc(ctx)
}

func (m *mockDockerSDK) Close() error {
	m.closeCalls++
	return m.closeFunc()
}

type mockReadCloser struct {
	io.Reader
	closeFunc func() error
}

func (m mockReadCloser) Close() error {
	return m.closeFunc()
}

func TestNewClient(t *testing.T) {
	t.Skip("Requires refactoring to allow mocking Docker client creation")
}

func TestRunningContainers(t *testing.T) {
	testCases := []struct {
		name             string
		mockContainers   []types.Container
		mockError        error
		expectedError    bool
		expectedContains []container.Container
		labelEnabled     bool
	}{
		{
			name: "successful container retrieval",
			mockContainers: []types.Container{
				{ID: "container1", Names: []string{"/test-container-1"}},
				{ID: "container2", Names: []string{"/test-container-2"}},
			},
			mockError: nil,
			expectedContains: []container.Container{
				{ID: "container1", Name: "test-container-1"},
				{ID: "container2", Name: "test-container-2"},
			},
			expectedError: false,
			labelEnabled:  false,
		},
		{
			name:             "error getting containers",
			mockContainers:   nil,
			mockError:        errors.New("docker error"),
			expectedContains: nil,
			expectedError:    true,
			labelEnabled:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSDK := &mockDockerSDK{
				containerListFunc: func(ctx context.Context, options types.ContainerListOptions) ([]types.Container, error) {
					return tc.mockContainers, tc.mockError
				},
				closeFunc: func() error {
					return nil
				},
			}

			client := &container.Client{
				SDK:  mockSDK,
				Opts: &container.ClientOptions{LabelEnabled: tc.labelEnabled},
			}

			containers, err := client.RunningContainers(context.Background())

			if tc.expectedError && err == nil {
				t.Error("expected error, but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if mockSDK.containerListCalls != 1 {
				t.Errorf("ContainerList should be called once, but was called %d times", mockSDK.containerListCalls)
			}

			if !tc.expectedError {
				if len(containers) != len(tc.expectedContains) {
					t.Errorf("expected %d containers, got %d", len(tc.expectedContains), len(containers))
				}

				for i, expected := range tc.expectedContains {
					if i >= len(containers) {
						t.Errorf("missing expected container %d", i)
						continue
					}
					actual := containers[i]
					if actual.ID != expected.ID || actual.Name != expected.Name {
						t.Errorf("container %d mismatch: expected %+v, got %+v", i, expected, actual)
					}
				}
			}
		})
	}
}

func TestContainerLogs(t *testing.T) {
	testCases := []struct {
		name           string
		containerID    string
		since          string
		tail           int
		mockLogs       string
		mockError      error
		readError      error
		expectedError  bool
		expectedOutput string
		ctx            context.Context
	}{
		{
			name:           "successful log retrieval",
			containerID:    "container1",
			since:          "1h",
			tail:           100,
			mockLogs:       "log line 1\nlog line 2",
			mockError:      nil,
			readError:      nil,
			expectedError:  false,
			expectedOutput: "log line 1\nlog line 2",
			ctx:            context.Background(),
		},
		{
			name:           "error getting logs",
			containerID:    "container1",
			since:          "1h",
			tail:           100,
			mockLogs:       "",
			mockError:      errors.New("docker error"),
			readError:      nil,
			expectedError:  true,
			expectedOutput: "",
			ctx:            context.Background(),
		},
		{
			name:           "error reading logs",
			containerID:    "container1",
			since:          "1h",
			tail:           100,
			mockLogs:       "log data",
			mockError:      nil,
			readError:      errors.New("read error"),
			expectedError:  true,
			expectedOutput: "",
			ctx:            context.Background(),
		},
		{
			name:           "empty container ID",
			containerID:    "",
			since:          "1h",
			tail:           100,
			mockLogs:       "",
			mockError:      nil,
			readError:      nil,
			expectedError:  true,
			expectedOutput: "",
			ctx:            context.Background(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSDK := &mockDockerSDK{
				containerLogsFunc: func(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
					if tc.mockError != nil {
						return nil, tc.mockError
					}
					return &mockReadCloser{
						Reader: strings.NewReader(tc.mockLogs),
						closeFunc: func() error {
							return nil
						},
					}, nil
				},
				closeFunc: func() error {
					return nil
				},
			}

			if tc.readError != nil {
				mockSDK.containerLogsFunc = func(ctx context.Context, container string, options types.ContainerLogsOptions) (io.ReadCloser, error) {
					return &mockReadCloser{
						Reader: &errorReader{err: tc.readError},
						closeFunc: func() error {
							return nil
						},
					}, nil
				}
			}

			client := &container.Client{
				SDK:  mockSDK,
				Opts: &container.ClientOptions{},
			}

			if tc.ctx == nil {
				defer func() {
					if r := recover(); r == nil {
						t.Error("expected panic with nil context, but none occurred")
					}
				}()
			}

			logs, err := client.ContainerLogs(tc.ctx, tc.containerID, tc.since, tc.tail)

			if tc.ctx == nil {
				return
			}

			if tc.expectedError && err == nil {
				t.Error("expected error, but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tc.containerID != "" {
				if mockSDK.containerLogsCalls != 1 {
					t.Errorf("ContainerLogs should be called once, but was called %d times", mockSDK.containerLogsCalls)
				}
			}

			if !tc.expectedError {
				if string(logs) != tc.expectedOutput {
					t.Errorf("expected logs '%s', got '%s'", tc.expectedOutput, string(logs))
				}
			}
		})
	}
}

type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

func TestClose(t *testing.T) {
	testCases := []struct {
		name          string
		mockError     error
		expectedError bool
		sdkIsNil      bool
	}{
		{
			name:          "successful close",
			mockError:     nil,
			expectedError: false,
			sdkIsNil:      false,
		},
		{
			name:          "error on close",
			mockError:     errors.New("close error"),
			expectedError: true,
			sdkIsNil:      false,
		},
		{
			name:          "sdk is nil",
			mockError:     nil,
			expectedError: false,
			sdkIsNil:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var client *container.Client

			if tc.sdkIsNil {
				client = &container.Client{
					SDK:  nil,
					Opts: &container.ClientOptions{},
				}
			} else {
				mockSDK := &mockDockerSDK{
					closeFunc: func() error {
						return tc.mockError
					},
				}

				client = &container.Client{
					SDK:  mockSDK,
					Opts: &container.ClientOptions{},
				}
			}

			err := client.Close()

			if tc.expectedError && err == nil {
				t.Error("expected error, but got none")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tc.sdkIsNil {
				mockSDK := client.SDK.(*mockDockerSDK)
				if mockSDK.closeCalls != 1 {
					t.Errorf("Close should be called once, but was called %d times", mockSDK.closeCalls)
				}
			}
		})
	}
}
