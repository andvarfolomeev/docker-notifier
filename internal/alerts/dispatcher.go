package alerts

import (
	"context"
	"log/slog"

	"github.com/andvarfolomeev/docker-notifier/internal/telegram"
	"github.com/andvarfolomeev/docker-notifier/internal/watcher"
)

func RunDispatcher(ctx context.Context, ch <-chan *watcher.MatchedLog, telegramClient *telegram.Client, log *slog.Logger) {
	for {
		select {
		case match, ok := <-ch:
			if !ok {
				return
			}

			message := PrepareMessage(match)
			if err := telegramClient.SendMessage(ctx, message); err != nil {
				log.Error("Failed to send message", "err", err)
			}

		case <-ctx.Done():
			log.Info("Context canceled, stopping dispatcher loop", "err", ctx.Err())
			return
		}
	}
}
