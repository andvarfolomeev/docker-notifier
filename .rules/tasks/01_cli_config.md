# 🏗 CLI and Configuration

## Goal
Implement a CLI interface with flags to configure docker-notifier behavior.

## Steps
1. Use [spf13/pflag](https://github.com/spf13/pflag) or the standard `flag` package.
2. Support the following flags:
    - `--interval`
    - `--label-enable`
    - `--telegram-token`
    - `--telegram-chat-id`
    - `--error-pattern` (repeatable)
    - `--debug`
    - `--cleanup`
3. Validate required flags (`telegram-token`, `telegram-chat-id`).
4. Display help when needed.

## Example
```bash
dockernotify --interval 30 --telegram-token "xxx" --telegram-chat-id "-xxx"
```