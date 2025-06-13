package alerts

import (
	"context"
	"log/slog"
	"time"

	"github.com/andvarfolomeev/docker-notifier/internal/telegram"
	"github.com/andvarfolomeev/docker-notifier/internal/watcher"
)

const timeout = 2 * time.Second

func RunDispatcher(ctx context.Context, ch <-chan *watcher.MatchedLog, telegramClient *telegram.Client, log *slog.Logger) {
	for {
		select {
		case match, ok := <-ch:
			if !ok {
				return
			}

			message := PrepareMessage(match)
			sendCtx, cancel := context.WithTimeout(ctx, timeout)
			err := telegramClient.SendMessage(sendCtx, message)
			cancel()
			if err != nil {
				log.Error("Failed to send message", "err", err)
			}

		case <-ctx.Done():
			log.Info("Context canceled, stopping dispatcher loop", "err", ctx.Err())
			return
		}
	}
}
