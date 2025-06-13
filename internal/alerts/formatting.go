package alerts

import (
	"fmt"
	"strings"

	"github.com/andvarfolomeev/docker-notifier/internal/watcher"
)

func PrepareMessage(match *watcher.MatchedLog) string {
	errorLine := match.Line.Content
	if len(errorLine) > 100 {
		errorLine = errorLine[:100]
	}

	messageLines := []string{
		fmt.Sprintf("🚨 Error detected!"),
		fmt.Sprintf("Container ID = %s; Container name = %s", match.Container.ID, match.Container.Name),
		fmt.Sprintf("Line: \"%s\"", errorLine),
	}
	message := strings.Join(messageLines, "\n")

	return message
}
