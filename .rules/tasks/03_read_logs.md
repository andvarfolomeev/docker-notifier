# 📜 Reading Container Logs

## Goal
Implement reading logs from containers.

## Steps
1. Use Docker API to retrieve logs.
2. On first run, read the **last N log lines** (e.g., 100 lines).
3. On subsequent iterations, read **only new lines** (incrementally).
4. Data structure: map `containerID -> lastOffset` (in memory).