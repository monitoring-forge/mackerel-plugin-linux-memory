VERSION=0.0.7
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
all: mackerel-plugin-linux-memory

.PHONY: mackerel-plugin-linux-memory linux check lint

mackerel-plugin-linux-memory: *.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-memory

linux: *.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-memory

check:
	go test -v ./...

lint:
	golangci-lint run ./...