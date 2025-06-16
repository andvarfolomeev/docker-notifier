# Docker Notifier 🐳

Monitor Docker container logs and get instant Telegram notifications when errors occur.

## Overview

Docker Notifier watches your container logs for error patterns and sends notifications to Telegram when matches are found. Easily integrate with your existing Docker setup to get alerts about errors in real-time.

## Features

- 🔍 Monitor container logs for custom error patterns
- 🏷️ Filter containers by labels
- 📱 Send notifications to Telegram
- ⏱️ Configurable polling interval
- 🛠️ Simple integration with docker-compose

## Quick Start

### Using docker-compose

```yaml
version: "3"

services:
  docker-notifier:
    image: ghcr.io/andvarfolomeev/docker-notifier:main
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    restart: unless-stopped
    environment:
      - TELEGRAM_TOKEN=your_telegram_bot_token
      - TELEGRAM_CHAT_ID=your_telegram_chat_id
    command: >
      /app/docker-notifier
        --interval 5
        --label-enable
        --telegram-token "${TELEGRAM_TOKEN}"
        --telegram-chat-id "${TELEGRAM_CHAT_ID}"
        --error-pattern "ERROR"
        --error-pattern "FATAL"
        --error-pattern "Exception"

  # Example service to monitor
  example-service:
    image: alpine:latest
    labels:
      - "com.andvarfolomeev.dockernotify.enable=true"
    command: >
      sh -c "while true; do echo 'Normal log line'; sleep 5; echo 'ERROR: This is an error message'; sleep 10; done"
```

### Using Docker run

```bash
docker run -d \
  --name docker-notifier \
  -v /var/run/docker.sock:/var/run/docker.sock \
  ghcr.io/andvarfolomeev/docker-notifier:main \
  /app/docker-notifier \
    --telegram-token "your_telegram_bot_token" \
    --telegram-chat-id "your_telegram_chat_id" \
    --error-pattern "ERROR" \
    --error-pattern "FATAL"
```

## Configuration

### Command-line Arguments

| Argument | Description | Default |
|----------|-------------|---------|
| `--interval` | Log polling interval in seconds | 5 |
| `--label-enable` | Enable label filter (only monitor containers with the label) | false |
| `--telegram-token` | Telegram Bot API token (required) | - |
| `--telegram-chat-id` | Target Telegram chat ID (required) | - |
| `--error-pattern` | Regex pattern for matching error lines (can be used multiple times) | "ERROR" |
| `--debug` | Enable debug logging | false |
| `--help` | Display help information | - |

### Container Labels

When `--label-enable` is set, Docker Notifier will only monitor containers with this label:

```
com.andvarfolomeev.dockernotify.enable=true
```

## Setup Telegram Bot

1. Create a new bot via [@BotFather](https://t.me/botfather) on Telegram
2. Get your bot token
3. Start a conversation with your bot
4. Get your chat ID (using [@userinfobot](https://t.me/userinfobot) or other methods)
5. Use these values for `--telegram-token` and `--telegram-chat-id`

## License

MIT
