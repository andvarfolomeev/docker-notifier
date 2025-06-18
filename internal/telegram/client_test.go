package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type transportWithURLOverride struct {
	base      http.RoundTripper
	serverURL string
}

func (t *transportWithURLOverride) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.Contains(req.URL.String(), "api.telegram.org") {
		newURL := fmt.Sprintf("%s/bot%s/sendMessage", t.serverURL, strings.Split(req.URL.Path, "/bot")[1])
		newReq := req.Clone(req.Context())
		parsedURL, err := req.URL.Parse(newURL)
		if err != nil {
			return nil, err
		}
		newReq.URL = parsedURL
		newReq.Host = parsedURL.Host
		return t.base.RoundTrip(newReq)
	}
	return t.base.RoundTrip(req)
}

func TestNew(t *testing.T) {
	token := "test-token"
	chatID := "test-chat-id"
	client := &http.Client{}

	telegramClient := New(token, chatID, client)

	if telegramClient.token != token {
		t.Errorf("expected token %s, got %s", token, telegramClient.token)
	}

	if telegramClient.chatID != chatID {
		t.Errorf("expected chatID %s, got %s", chatID, telegramClient.chatID)
	}

	if telegramClient.client != client {
		t.Errorf("expected client to be the same instance")
	}
}

func TestSendMessage(t *testing.T) {
	var receivedRequest *http.Request
	var requestBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedRequest = r
		var err error
		requestBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Failed to read request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	token := "test-token"
	chatID := "test-chat-id"
	message := "test message"
	customClient := &Client{
		token:  token,
		chatID: chatID,
		client: &http.Client{
			Transport: &transportWithURLOverride{
				base:      server.Client().Transport,
				serverURL: server.URL,
			},
		},
	}

	err := customClient.SendMessage(context.Background(), message)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if receivedRequest.Method != http.MethodPost {
		t.Errorf("expected POST request, got %s", receivedRequest.Method)
	}

	if contentType := receivedRequest.Header.Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var msgReq MessageRequest
	if err := json.Unmarshal(requestBody, &msgReq); err != nil {
		t.Fatalf("Failed to unmarshal request body: %v", err)
	}

	if msgReq.ChatID != chatID {
		t.Errorf("expected ChatID %s, got %s", chatID, msgReq.ChatID)
	}

	if msgReq.Text != message {
		t.Errorf("expected Text %s, got %s", message, msgReq.Text)
	}
}

func TestSendMessage_Error(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse int
		expectedError  bool
	}{
		{
			name:           "Success",
			serverResponse: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "Telegram API Error",
			serverResponse: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.serverResponse)
				if tc.serverResponse == http.StatusOK {
					w.Write([]byte(`{"ok":true}`))
				} else {
					w.Write([]byte(`{"ok":false}`))
				}
			}))
			defer server.Close()

			token := "test-token"
			chatID := "test-chat-id"
			message := "test message"
			customClient := &Client{
				token:  token,
				chatID: chatID,
				client: &http.Client{
					Transport: &transportWithURLOverride{
						base:      server.Client().Transport,
						serverURL: server.URL,
					},
				},
			}

			err := customClient.SendMessage(context.Background(), message)

			if tc.expectedError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tc.expectedError && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
	}
}
