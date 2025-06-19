package container_test

import (
	"fmt"
	"testing"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
	"github.com/andvarfolomeev/docker-notifier/internal/docker"
)

func TestRunningContainerFilters(t *testing.T) {
	testCases := []struct {
		name          string
		labelEnabled  bool
		expectedCount int
	}{
		{
			name:          "status filter only",
			labelEnabled:  false,
			expectedCount: 1,
		},
		{
			name:          "status and label filters",
			labelEnabled:  true,
			expectedCount: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := &container.ClientOptions{
				LabelEnabled: tc.labelEnabled,
			}

			filterArgs := container.RunningContainerFilters(opts)

			expectedFilter := docker.NewFilter()
			expectedFilter.Add("status", "running")

			filterCount := 1

			if tc.labelEnabled {
				expectedLabelValue := fmt.Sprintf("%s=%s", container.LabelEnableKey, container.LabelEnableValue)
				expectedFilter.Add("label", expectedLabelValue)
				filterCount++
			}

			expected, _ := expectedFilter.Encode()
			actual, _ := filterArgs.Encode()

			if expected != actual {
				t.Errorf("Filter mismatch: expected %s, got %s", expected, actual)
			}

			if filterArgs.Len() != tc.expectedCount {
				t.Errorf("expected %d filters, got %d", tc.expectedCount, filterArgs.Len())
			}
		})
	}
}

func TestContainerLogsOptions(t *testing.T) {
	testCases := []struct {
		name               string
		since              string
		tail               int
		expectedSince      string
		expectedTail       string
		expectedShowStdout bool
		expectedShowStderr bool
		expectedTimestamps bool
		expectedFollow     bool
	}{
		{
			name:               "no parameters",
			since:              "",
			tail:               0,
			expectedSince:      "",
			expectedTail:       "",
			expectedShowStdout: true,
			expectedShowStderr: true,
			expectedTimestamps: true,
			expectedFollow:     false,
		},
		{
			name:               "with since parameter",
			since:              "1h",
			tail:               0,
			expectedSince:      "1h",
			expectedTail:       "",
			expectedShowStdout: true,
			expectedShowStderr: true,
			expectedTimestamps: true,
			expectedFollow:     false,
		},
		{
			name:               "with tail parameter",
			since:              "",
			tail:               100,
			expectedSince:      "",
			expectedTail:       "100",
			expectedShowStdout: true,
			expectedShowStderr: true,
			expectedTimestamps: true,
			expectedFollow:     false,
		},
		{
			name:               "with since and tail parameters",
			since:              "2h",
			tail:               50,
			expectedSince:      "2h",
			expectedTail:       "50",
			expectedShowStdout: true,
			expectedShowStderr: true,
			expectedTimestamps: true,
			expectedFollow:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := container.ContainerLogsOptions(tc.since, tc.tail)

			if opts.Since != tc.expectedSince {
				t.Errorf("Since field: expected '%s', got '%s'", tc.expectedSince, opts.Since)
			}
			if opts.Tail != tc.expectedTail {
				t.Errorf("Tail field: expected '%s', got '%s'", tc.expectedTail, opts.Tail)
			}
			if opts.Stdout != tc.expectedShowStdout {
				t.Errorf("Stdout field: expected %v, got %v", tc.expectedShowStdout, opts.Stdout)
			}
			if opts.Stderr != tc.expectedShowStderr {
				t.Errorf("Stderr field: expected %v, got %v", tc.expectedShowStderr, opts.Stderr)
			}
			if opts.Timestamp != tc.expectedTimestamps {
				t.Errorf("Timestamp field: expected %v, got %v", tc.expectedTimestamps, opts.Timestamp)
			}
			if opts.Follow != tc.expectedFollow {
				t.Errorf("Follow field: expected %v, got %v", tc.expectedFollow, opts.Follow)
			}
		})
	}
}
