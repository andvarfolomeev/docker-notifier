package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const telegramAPIURL = "https://api.telegram.org/bot%s/sendMessage"

type Client struct {
	token  string
	chatID string
	client *http.Client
}

func New(token, chatID string, client *http.Client) *Client {
	c := &Client{
		token:  token,
		chatID: chatID,
		client: client,
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

	return nil
}
