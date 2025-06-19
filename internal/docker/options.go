package docker

type ContainerListOptions struct {
	Size    bool
	All     bool
	Latest  bool
	Since   string
	Before  string
	Limit   int
	Filters Filters
}

type ContainerLogsOptions struct {
	Follow    bool
	Stdout    bool
	Stderr    bool
	Since     string
	Until     string
	Timestamp bool
	Tail      string
}
