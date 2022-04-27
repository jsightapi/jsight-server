.PHONY: all
all: fmt lint build

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover ./...

.PHONY: build
build:
	go build -o jsight-server .
