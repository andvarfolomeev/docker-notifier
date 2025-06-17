package alerts_test

import (
	"testing"

	"github.com/andvarfolomeev/docker-notifier/internal/alerts"
	"github.com/andvarfolomeev/docker-notifier/internal/container"
	"github.com/andvarfolomeev/docker-notifier/internal/logfilter"
	"github.com/andvarfolomeev/docker-notifier/internal/watcher"
)

func TestPrepareMessage(t *testing.T) {
	testCases := []struct {
		name     string
		match    *watcher.MatchedLog
		expected string
	}{
		{
			name: "regular error message",
			match: &watcher.MatchedLog{
				Container: container.Container{
					ID:   "abc123",
					Name: "test-container",
				},
				Line: &logfilter.MatchedLine{
					Content: []byte("Error: connection refused"),
				},
			},
			expected: "🚨 Error detected!\nContainer ID = abc123; Container name = test-container\nLine: \"Error: connection refused\"",
		},
		{
			name: "very long error message gets truncated",
			match: &watcher.MatchedLog{
				Container: container.Container{
					ID:   "def456",
					Name: "long-error-container",
				},
				Line: &logfilter.MatchedLine{
					Content: []byte("Error: very long error message that exceeds 100 characters and should be truncated by the formatting function to keep messages concise and readable"),
				},
			},
			expected: "🚨 Error detected!\nContainer ID = def456; Container name = long-error-container\nLine: \"Error: very long error message that exceeds 100 characters and should be truncated by the formatting\"",
		},
		{
			name: "empty error message",
			match: &watcher.MatchedLog{
				Container: container.Container{
					ID:   "ghi789",
					Name: "empty-error-container",
				},
				Line: &logfilter.MatchedLine{
					Content: []byte(""),
				},
			},
			expected: "🚨 Error detected!\nContainer ID = ghi789; Container name = empty-error-container\nLine: \"\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			message := alerts.PrepareMessage(tc.match)
			if message != tc.expected {
				t.Errorf("expected message: %q, got: %q", tc.expected, message)
			}
		})
	}
}
