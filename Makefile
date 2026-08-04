VERSION=0.0.7
GITCOMMIT?=$(shell git describe --dirty --always 2>/dev/null)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"
all: mackerel-plugin-linux-memory

.PHONY: mackerel-plugin-linux-process-status

mackerel-plugin-linux-memory: *.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-memory

linux: *.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-memory

check:
	go test -v ./...

lint:
	golangci-lint run ./...