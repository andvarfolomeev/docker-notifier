# 🔄 Handling Container Restarts

## Goal
Handle container restarts.

## Steps
1. Listen to Docker events (Container Started / Stopped).
2. On container restart, reset offset and start reading logs from the beginning.
3. Implement offset storage in memory.

4. Implement `--cleanup`:
    - Reset in-memory offsets.