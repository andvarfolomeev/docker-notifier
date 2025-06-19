package docker

import (
	"fmt"
	"time"
)

func parseTimestamp(s string) (int64, error) {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return 0, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	return t.Unix(), nil
}
