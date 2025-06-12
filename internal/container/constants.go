package container

import "time"

const (
	defaultDockerSocket    = "unix:///var/run/docker.sock"
	labelEnableKey         = "com.andvarfolomeev.dockernotify.enable"
	labelEnableValue       = "true"
	defaultInitialLogLines = 100
	pingTimeout            = 5 * time.Second
	logBufferSize          = 64 * 1024
	shortIDLen             = 12
)
