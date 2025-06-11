# 🔍 Matching Error Patterns

## Goal
Match log lines against error patterns.

## Steps
1. Use Go's `regexp` package.
2. Support multiple `--error-pattern` flags (repeatable).
3. If a match is found, send an alert to Telegram.
4. Support OR logic for patterns.