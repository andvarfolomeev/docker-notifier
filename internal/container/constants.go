package container

import "time"

const (
	LabelEnableKey   = "com.andvarfolomeev.dockernotify.enable"
	LabelEnableValue = "true"
	PingTimeout      = 5 * time.Second
	ShortIDLen       = 12
)
