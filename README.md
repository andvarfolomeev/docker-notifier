# 🐳 Docker Notifier

A simple service that connects to `docker.sock`, watches logs of selected containers, and sends error messages to Telegram.

## Features

- Monitors Docker container logs for error patterns
- Sends alerts to Telegram when errors are detected
- Configurable polling interval
- Container filtering via Docker labels
- Multiple error pattern matching via regex
- Small, self-contained binary with no dependencies

## Usage

```bash
dockernotify                    \
    --interval 30               \
    --label-enable              \
    --telegram-token "xxx"      \
    --telegram-chat-id "-xxx"   \
    --error-pattern "ERROR"     \
    --error-pattern "FATAL"     \
    --error-pattern "Exception" \
    --debug                     \
    --cleanup
```

## Supported Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--interval <seconds>` | Log polling interval | `5` |
| `--label-enable` | Enable label filter: `com.andvarfolomeev.dockernotify.enable=true` | *flag* |
| `--telegram-token <token>` | Telegram Bot API token | — |
| `--telegram-chat-id <chat_id>` | Target chat ID | — |
| `--error-pattern <pattern>` | Regex pattern for matching error lines (can be used multiple times) | `"ERROR"` |
| `--debug` | Enable debug logging | *flag* |
| `--cleanup` | (Optional) Clear saved log offsets | *flag* |

## Docker Compose Example

```yaml
services:
  dockernotify:
    image: your-dockernotify-image
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    command: >
      --interval 30
      --label-enable
      --telegram-token "xxx"
      --telegram-chat-id "-xxx"
      --error-pattern "ERROR"
      --error-pattern "FATAL"
```

## Building from Source

```bash
go build -o bin/dockernotify cmd/docker-notifier/main.go
```

## License

[MIT License](LICENSE)