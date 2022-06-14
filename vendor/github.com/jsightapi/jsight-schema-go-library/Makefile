.PHONY: all
all: generate fmt lint test

.PHONY: generate
generate:
	go generate $$(go list ./... | grep -v vendor)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -cover ./...

.PHONY: bench
bench:
	go test -run XXXX -bench . -benchmem ./...