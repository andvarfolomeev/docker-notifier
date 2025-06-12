# 🐳 Docker Notifier

A simple service that connects to `docker.sock`, watches logs of selected containers, and sends error messages to Telegram. The service implements tasks 02 (Docker connection) and 03 (Reading logs).

## ⚠️ Docker Socket Permissions

When running this service, you may encounter permission issues with the Docker socket (`/var/run/docker.sock`). This happens because the Docker socket is owned by the `docker` group on the host, and the container needs appropriate permissions to access it.

There are several ways to solve this:

1. **Use the start script** (recommended):
   ```bash
   ./scripts/start.sh
   ```
   This script automatically detects the Docker group ID and starts the container with the correct permissions.

2. **Set group ID manually**:
   ```bash
   # Find Docker group ID
   DOCKER_GID=$(stat -c '%g' /var/run/docker.sock)
   
   # Run with the correct group
   docker-compose up -d -e GID=${DOCKER_GID}
   ```

3. **Run as root** (not recommended for production):
   This is less secure but simpler for testing:
   ```bash
   docker-compose up -d
   ```

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
    # Fix permission issues with Docker socket
    user: "${UID:-0}:${GID:-0}"
    command: >
      --interval 30
      --label-enable
      --telegram-token "xxx"
      --telegram-chat-id "-xxx"
      --error-pattern "ERROR"
      --error-pattern "FATAL"
```

The `user` parameter ensures that the container has the correct permissions to access the Docker socket. The environment variables `UID` and `GID` should be set to match the user and group IDs of the Docker socket on your host.

## Building from Source

```bash
go build -o bin/dockernotify ./cmd/docker-notifier
```

## Quick Start

1. Clone the repository
2. Build the Docker image:
   ```bash
   make docker-build
   ```
3. Set your Telegram token and chat ID in `docker-compose.yml`
4. Run the service with the start script:
   ```bash
   ./scripts/start.sh
   ```
5. Test the service by running a container with the required label:
   ```bash
   docker run -d --name test-container \
     --label com.andvarfolomeev.dockernotify.enable=true \
     alpine sh -c "while true; do echo 'ERROR: Test error message'; sleep 10; done"
   ```

## Implementation Details

### Docker API Connection
- Connects to Docker via Unix socket at `/var/run/docker.sock`
- Uses the official Docker Go client library
- Filters running containers with label `com.andvarfolomeev.dockernotify.enable=true` when `--label-enable` is specified

### Log Reading
- On first run: reads the last 100 log lines from each container
- On subsequent polling intervals: reads only new log lines incrementally
- Maintains in-memory offsets to track log position for each container
- Parses log timestamps to enable incremental reading
- Scans logs for configurable error patterns

## License

[MIT License](LICENSE)