package watcher

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"sync"
	"time"

	"github.com/andvarfolomeev/docker-notifier/internal/container"
	"github.com/andvarfolomeev/docker-notifier/internal/logfilter"
)

const maxConsecutiveFailures = 5

type MatchedLog struct {
	Container container.Container
	Line      *logfilter.MatchedLine
}

type Watcher struct {
	client   ContainerClient
	interval time.Duration
	patterns []*regexp.Regexp
	C        chan *MatchedLog

	mu      sync.RWMutex
	offsets map[string]string
}

type WatcherOptions struct {
	Interval      time.Duration
	ErrorPatterns []string
}

func New(
	client ContainerClient,
	opts *WatcherOptions,
) (*Watcher, error) {
	patterns, err := compileErrorPatterns(opts.ErrorPatterns)
	if err != nil {
		return nil, err
	}

	offsets := make(map[string]string)

	c := make(chan *MatchedLog)

	w := &Watcher{
		client:   client,
		interval: opts.Interval,
		patterns: patterns,
		offsets:  offsets,
		C:        c,
	}

	return w, nil
}

func (w *Watcher) Start(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)

	go func() {
		defer cancel()
		w.start(ctx)
	}()
}

func (w *Watcher) start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	consecutiveFailures := 0

	for {
		select {
		case <-ticker.C:
			if err := w.checkContainers(ctx); err != nil {
				consecutiveFailures++
			} else {
				consecutiveFailures = 0
			}

			if consecutiveFailures >= maxConsecutiveFailures {
				slog.Error("Too many consecutive failures")
				return
			}
		case <-ctx.Done():
			return
		}
	}

}

func (w *Watcher) checkContainers(ctx context.Context) error {
	containers, err := w.client.RunningContainers(ctx)
	if err != nil {
		return fmt.Errorf("Failed to list containers: %w", err)
	}

	for _, container := range containers {
		if err := w.processContainerLogs(ctx, container); err != nil {
			slog.Error("Failed to process container logs", "err", err)
		}

	}

	return nil
}

func (w *Watcher) processContainerLogs(ctx context.Context, container container.Container) error {
	w.mu.RLock()
	since, ok := w.offsets[container.ID]
	w.mu.RUnlock()

	if !ok {
		w.mu.Lock()
		w.offsets[container.ID] = nowStrSince()
		w.mu.Unlock()
		slog.Debug("First time seeing container", "containerID", container.ID)
		return nil
	}

	lines, err := w.client.ContainerLogs(ctx, container.ID, since, 0)
	if err != nil {
		return fmt.Errorf("Failed to to get logs for container %s: %w", container.ID, err)
	}

	matchedLogs, err := logfilter.FindMatchedLines(w.patterns, lines)
	if err != nil {
		return fmt.Errorf("Failed to process logs for container %s: %w", container.ID, err)
	}

	sinceTime, err := parseStrSince(since)
	if err != nil {
		return fmt.Errorf("Failed to process logs for container %s: %w", container.ID, err)
	}

	for _, matchedLog := range matchedLogs {
		lastSinceTime, err := parseStrSince(string(matchedLog.Timestamp))
		if err != nil {
			return fmt.Errorf("Failed to process logs for container %s: %w", container.ID, err)
		}

		if !lastSinceTime.After(sinceTime) {
			continue
		}

		m := &MatchedLog{
			Container: container,
			Line:      matchedLog,
		}

		select {
		case w.C <- m:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	if len(matchedLogs) > 0 {
		newOffset := string(matchedLogs[len(matchedLogs)-1].Timestamp)

		w.mu.Lock()
		w.offsets[container.ID] = newOffset
		w.mu.Unlock()
	}

	return nil
}

func (w *Watcher) Cleanup() {
	close(w.C)
}

func compileErrorPatterns(errorPatterns []string) ([]*regexp.Regexp, error) {
	patterns := make([]*regexp.Regexp, 0, len(errorPatterns))
	for _, pattern := range errorPatterns {
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid error pattern '%s': %w", pattern, err)
		}
		patterns = append(patterns, re)
	}
	return patterns, nil
}

func nowStrSince() string {
	return time.Now().Format(time.RFC3339Nano)
}

func parseStrSince(since string) (time.Time, error) {
	ts, err := time.Parse(time.RFC3339Nano, since)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse since %s: %w", since, err)
	}
	return ts, nil
}
