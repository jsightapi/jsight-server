.PHONY: all
all: fmt lint build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: build
build:
	go build -o jsight-server .
