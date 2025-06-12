package telegram

import (
	"log/slog"
	"net/http"
)

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.client = httpClient
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) {
		c.log = logger
	}
}
