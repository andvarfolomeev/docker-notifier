module github.com/andvarfolomeev/docker-notifier

go 1.23.0

toolchain go1.23.10

require github.com/spf13/pflag v1.0.6

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/docker/distribution v0.0.0-00010101000000-000000000000 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/moby/term v0.5.2 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/time v0.12.0 // indirect
)

require (
	github.com/docker/docker v24.0.9+incompatible
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	gotest.tools/v3 v3.5.2 // indirect
)

replace github.com/docker/distribution => github.com/docker/distribution v2.8.1+incompatible
