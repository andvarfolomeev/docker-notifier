# 🐞 Logging and Debugging

## Goal
Add support for debug logging.

## Steps
1. Use the `log` package or external libraries (`zerolog`, `zap`).
2. Add the `--debug` flag:
    - Enables verbose output.
3. Log the following:
    - Connection to docker.sock.
    - Discovered containers.
    - New log lines read.
    - Matches with error patterns.
    - Messages sent to Telegram.