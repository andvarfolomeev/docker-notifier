package watcher

import "time"

const (
	defaultInterval        = 5 * time.Second
	maxConsecutiveFailures = 5
)
