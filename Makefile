VERSION=0.0.7
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
all: mackerel-plugin-linux-memory

.PHONY: mackerel-plugin-linux-process-status

mackerel-plugin-linux-memory: main.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-memory

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-memory

fmt:
	go fmt ./...

check:
	go test -v ./...

lint:
	golangci-lint run ./...