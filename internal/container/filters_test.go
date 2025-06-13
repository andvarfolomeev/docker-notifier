package container_test

import (
	"fmt"
	"testing"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
	"github.com/docker/docker/api/types/filters"
)

func TestRunningContainerFilters(t *testing.T) {
	testCases := []struct {
		name          string
		labelEnabled  bool
		expectedArgs  filters.Args
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

			// Check that status filter is always present
			statusValues := filterArgs.Get("status")
			if len(statusValues) != 1 || statusValues[0] != "running" {
				t.Errorf("expected filter status=running, got %v", statusValues)
			}

			// Check if label is added or not based on settings
			labelValues := filterArgs.Get("label")
			expectedLabelValue := fmt.Sprintf("%s=%s", container.LabelEnableKey, container.LabelEnableValue)

			if tc.labelEnabled {
				if len(labelValues) != 1 || labelValues[0] != expectedLabelValue {
					t.Errorf("expected filter label=%s, got %v", expectedLabelValue, labelValues)
				}
			} else {
				if len(labelValues) != 0 {
					t.Errorf("label filter not expected, but got %v", labelValues)
				}
			}

			// Check total number of filters
			if len(filterArgs.Keys()) != tc.expectedCount {
				t.Errorf("expected %d filters, got %d", tc.expectedCount, len(filterArgs.Keys()))
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

			// Check all field values
			if opts.Since != tc.expectedSince {
				t.Errorf("Since field: expected '%s', got '%s'", tc.expectedSince, opts.Since)
			}
			if opts.Tail != tc.expectedTail {
				t.Errorf("Tail field: expected '%s', got '%s'", tc.expectedTail, opts.Tail)
			}
			if opts.ShowStdout != tc.expectedShowStdout {
				t.Errorf("ShowStdout field: expected %v, got %v", tc.expectedShowStdout, opts.ShowStdout)
			}
			if opts.ShowStderr != tc.expectedShowStderr {
				t.Errorf("ShowStderr field: expected %v, got %v", tc.expectedShowStderr, opts.ShowStderr)
			}
			if opts.Timestamps != tc.expectedTimestamps {
				t.Errorf("Timestamps field: expected %v, got %v", tc.expectedTimestamps, opts.Timestamps)
			}
			if opts.Follow != tc.expectedFollow {
				t.Errorf("Follow field: expected %v, got %v", tc.expectedFollow, opts.Follow)
			}
		})
	}
}
