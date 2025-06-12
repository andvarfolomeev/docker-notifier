package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	telegramAPIURL = "https://api.telegram.org/bot%s/sendMessage"
	timeout        = 10 * time.Second
)

type Client struct {
	token  string
	chatID string
	log    *slog.Logger
	client *http.Client
}

func New(token, chatID string, log *slog.Logger, opts ...Option) *Client {
	if log == nil {
		log = slog.Default()
	}

	c := &Client{
		token:  token,
		chatID: chatID,
		log:    log,
		client: &http.Client{
			Timeout: timeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) SendMessage(ctx context.Context, message string) error {
	url := fmt.Sprintf(telegramAPIURL, c.token)

	requestBody, err := json.Marshal(MessageRequest{
		ChatID: c.chatID,
		Text:   message,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram API returned non-OK status: %s", resp.Status)
	}

	c.log.Debug("Successfully sent message to Telegram")
	return nil
}
