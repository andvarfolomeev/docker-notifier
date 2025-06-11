# 🐳 Connecting to Docker API

## Goal
Connect to the Docker API via `docker.sock`.

## Steps
1. Use [docker/docker/client](https://github.com/moby/moby/tree/master/client) or [docker/docker](https://pkg.go.dev/github.com/docker/docker).
2. Connect to `/var/run/docker.sock` (unix socket).
3. Retrieve the list of **running** containers.
4. Filter containers with label `com.andvarfolomeev.dockernotify.enable=true`.