package watcher

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
)

func TestNew(t *testing.T) {
	client := NewMockContainerClient()

	tests := []struct {
		name          string
		opts          *WatcherOptions
		expectedError bool
	}{
		{
			name: "valid options",
			opts: &WatcherOptions{
				Interval:      time.Second,
				ErrorPatterns: []string{"ERROR", "FATAL"},
			},
			expectedError: false,
		},
		{
			name: "invalid regex pattern",
			opts: &WatcherOptions{
				Interval:      time.Second,
				ErrorPatterns: []string{"[invalid"},
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := New(client, tt.opts)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if w != nil {
					t.Errorf("expected nil watcher, got %v", w)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if w == nil {
					t.Errorf("expected watcher, got nil")
				}
				if w.interval != tt.opts.Interval {
					t.Errorf("expected interval %v, got %v", tt.opts.Interval, w.interval)
				}
				if len(w.patterns) != len(tt.opts.ErrorPatterns) {
					t.Errorf("expected %d patterns, got %d", len(tt.opts.ErrorPatterns), len(w.patterns))
				}
			}
		})
	}
}

func TestWatcher_processContainerLogs(t *testing.T) {
	tests := []struct {
		name            string
		container       container.Container
		patterns        []string
		logs            string
		expectedMatches int
	}{
		{
			name: "match error logs",
			container: container.Container{
				ID:   "container1",
				Name: "test-container",
			},
			patterns: []string{"ERROR"},
			logs: "2023-03-15T12:02:00.000000000Z ERROR: Connection failed\n" +
				"2023-03-15T12:05:00.000000000Z ERROR: Database timeout",
			expectedMatches: 2,
		},
		{
			name: "match warning logs",
			container: container.Container{
				ID:   "container2",
				Name: "test-container",
			},
			patterns: []string{"WARNING"},
			logs: "2023-03-15T12:01:00.000000000Z Processing request\n" +
				"2023-03-15T12:04:00.000000000Z WARNING: High memory usage",
			expectedMatches: 1,
		},
		{
			name: "no matches",
			container: container.Container{
				ID:   "container3",
				Name: "test-container",
			},
			patterns: []string{"CRITICAL"},
			logs: "2023-03-15T12:00:00.000000000Z Container started\n" +
				"2023-03-15T12:07:00.000000000Z Processing complete",
			expectedMatches: 0,
		},
		{
			name: "multiple patterns",
			container: container.Container{
				ID:   "container4",
				Name: "test-container",
			},
			patterns: []string{"ERROR", "WARNING"},
			logs: "2023-03-15T12:02:00.000000000Z ERROR: Connection failed\n" +
				"2023-03-15T12:04:00.000000000Z WARNING: High memory usage\n" +
				"2023-03-15T12:05:00.000000000Z ERROR: Database timeout",
			expectedMatches: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewMockContainerClient()
			client.SetContainers([]container.Container{tt.container})
			client.SetLogs(tt.container.ID, []byte(tt.logs))

			patterns, err := compileErrorPatterns(tt.patterns)
			if err != nil {
				t.Fatalf("failed to compile patterns: %v", err)
			}

			watcher := &Watcher{
				client:   client,
				interval: time.Millisecond * 10,
				patterns: patterns,
				offsets:  make(map[string]string),
				C:        make(chan *MatchedLog, 10),
			}

			initTime := "2023-03-15T12:00:00.000000000Z"
			watcher.mu.Lock()
			watcher.offsets[tt.container.ID] = initTime
			watcher.mu.Unlock()

			err = watcher.processContainerLogs(context.Background(), tt.container)
			if err != nil {
				t.Fatalf("processContainerLogs failed: %v", err)
			}

			matches := 0
			for {
				select {
				case <-watcher.C:
					matches++
				case <-time.After(time.Millisecond * 50):
					goto done
				}
			}
		done:

			if matches != tt.expectedMatches {
				t.Errorf("expected %d matches, got %d", tt.expectedMatches, matches)
			}

			close(watcher.C)
		})
	}
}

func TestWatcher_checkContainers(t *testing.T) {
	client := NewMockContainerClient()

	containers := []container.Container{
		{ID: "container1", Name: "test-container1"},
		{ID: "container2", Name: "test-container2"},
	}

	client.SetContainers(containers)
	client.SetLogs("container1", []byte("2023-03-15T12:02:00.000000000Z ERROR: Connection failed\n2023-03-15T12:05:00.000000000Z ERROR: Database timeout"))
	client.SetLogs("container2", []byte("2023-03-15T12:04:00.000000000Z WARNING: High memory usage"))

	opts := &WatcherOptions{
		Interval:      time.Millisecond * 10,
		ErrorPatterns: []string{"ERROR"},
	}

	watcher, err := New(client, opts)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}

	err = watcher.checkContainers(context.Background())
	if err != nil {
		t.Fatalf("checkContainers failed: %v", err)
	}

	if client.ContainersCallCount() != 1 {
		t.Errorf("expected 1 container call, got %d", client.ContainersCallCount())
	}

	err = watcher.checkContainers(context.Background())
	if err != nil {
		t.Fatalf("checkContainers failed: %v", err)
	}

	if client.ContainersCallCount() != 2 {
		t.Errorf("expected 2 container calls, got %d", client.ContainersCallCount())
	}
}

func TestWatcher_Start(t *testing.T) {
	client := NewMockContainerClient()

	containers := []container.Container{
		{ID: "container1", Name: "test-container1"},
	}

	client.SetContainers(containers)
	client.SetLogs("container1", []byte("2023-03-15T12:02:00.000000000Z ERROR: Connection failed"))

	opts := &WatcherOptions{
		Interval:      time.Millisecond * 50,
		ErrorPatterns: []string{"ERROR"},
	}

	watcher, err := New(client, opts)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	watcher.Start(ctx)

	time.Sleep(time.Millisecond * 120)

	cancel()

	time.Sleep(time.Millisecond * 10)

	callCount := client.ContainersCallCount()
	if callCount < 2 {
		t.Errorf("expected at least 2 container checks, got %d", callCount)
	}
}

func TestHelperFunctions(t *testing.T) {
	t.Run("compileErrorPatterns", func(t *testing.T) {
		patterns := []string{"ERROR", "FATAL", "CRITICAL"}
		compiled, err := compileErrorPatterns(patterns)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(compiled) != len(patterns) {
			t.Errorf("expected %d patterns, got %d", len(patterns), len(compiled))
		}

		_, err = compileErrorPatterns([]string{"["})
		if err == nil {
			t.Error("expected error for invalid pattern, got nil")
		}
	})

	t.Run("nowStrSince", func(t *testing.T) {
		timeStr := nowStrSince()
		_, err := time.Parse(time.RFC3339Nano, timeStr)
		if err != nil {
			t.Errorf("failed to parse time: %v", err)
		}
	})

	t.Run("parseStrSince", func(t *testing.T) {
		validTime := "2023-03-15T12:00:00.000000000Z"
		parsed, err := parseStrSince(validTime)
		if err != nil {
			t.Errorf("failed to parse time: %v", err)
		}
		if parsed.Year() != 2023 {
			t.Errorf("expected year 2023, got %d", parsed.Year())
		}
		if parsed.Month() != time.March {
			t.Errorf("expected month March, got %s", parsed.Month())
		}
		if parsed.Day() != 15 {
			t.Errorf("expected day 15, got %d", parsed.Day())
		}

		_, err = parseStrSince("invalid-time")
		if err == nil {
			t.Error("expected error for invalid time, got nil")
		}
	})
}

func TestWatcherConsecutiveFailures(t *testing.T) {
	client := NewMockContainerClient()

	expectedErr := errors.New("container error")
	client.SetContainersError(expectedErr)

	opts := &WatcherOptions{
		Interval:      time.Millisecond * 10,
		ErrorPatterns: []string{"ERROR"},
	}

	watcher, err := New(client, opts)
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	watcher.start(ctx)

	callCount := client.ContainersCallCount()
	if callCount < maxConsecutiveFailures {
		t.Errorf("expected at least %d container checks, got %d", maxConsecutiveFailures, callCount)
	}
}
