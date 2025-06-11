# 🐳 docker notifier

## Goal

A simple service that connects to `docker.sock`, watches logs of selected containers, and sends error messages to Telegram.

---

## CLI Interface

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

---

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

---

## Behavior

✅ **Startup**
- Connect to `docker.sock`.
- Scan all running containers.
- For containers with label `com.andvarfolomeev.dockernotify.enable=true`, start reading logs.
- On first run: read **last N log lines** (e.g. 100 lines).

✅ **Polling**
- Every `--interval` seconds, read new logs (incremental).
- Match log lines against `--error-pattern`.
- If a match is found, send a message to Telegram.

✅ **Handling Restarts**
- If a container restarts, restart log reading from the beginning.

✅ **Offsets**
- Initially in-memory (no persistent storage required).
- Optionally support `--cleanup` to reset offsets.

✅ **Telegram Message Example**

```
🚨 [my-cool-service] Error detected!

Line: "2025-06-11 12:34:56 ERROR something bad happened"

Container: my-cool-service
```

---

## Architecture

- Written in **Go** (fast & small binary).
- Single binary — no dependencies.
- Runs as container with `docker.sock` mounted.

### Example docker-compose.yml

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

---

## Possible Future Improvements

- Persistent offset storage (BoltDB / simple file).
- Deduplication of repeated alerts.
- Markdown formatting of Telegram messages.
- Support for startup grace period (ignore errors on container startup).

---

## Summary

A small, fast, self-contained tool to monitor container logs for errors and send alerts to Telegram — without requiring complex monitoring stacks.

---
