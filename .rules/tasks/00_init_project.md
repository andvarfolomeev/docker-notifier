# 📦 Project Initialization

## Goal
Set up the Go project structure for developing docker-notifier.

## Steps
1. Initialize the go module:
    ```bash
    go mod init github.com/yourusername/docker-notifier
    ```

2. Create basic project directories:
    ```
    /cmd/docker-notifier/main.go
    /internal/config/
    /internal/docker/
    /internal/logger/
    /internal/telegram/
    /internal/watcher/
    /pkg/
    /README.md
    ```

3. Configure `.gitignore`:
    ```
    /bin/
    /vendor/
    *.log
    ```

4. Initialize Git repository:
    ```bash
    git init
    git commit -m "Initial project setup"
    ```