package container

import "time"

const (
	labelEnableKey   = "com.andvarfolomeev.dockernotify.enable"
	labelEnableValue = "true"
	pingTimeout      = 5 * time.Second
	shortIDLen       = 12
)
