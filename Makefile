#!/usr/bin/make -f

export CGO_ENABLED=0

PROJECT=github.com/previousnext/imaginator

# Builds the project.
build:
	gox -os='linux darwin' \
	    -arch='amd64' \
	    -output='bin/imaginator_{{.OS}}_{{.Arch}}' \
	    -ldflags='-extldflags "-static"' \
	    $(PROJECT)

# Run all lint checking with exit codes for CI.
lint:
	golint -set_exit_status `go list ./... | grep -v /vendor/`

# Run tests with coverage reporting.
test:
	go test -cover ./...

.PHONY: *
