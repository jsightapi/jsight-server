.PHONY: all
all: generate fmt lint test

.PHONY: deps
deps:
	go install github.com/vektra/mockery/v2@v2.12.3
	go install golang.org/x/tools/cmd/stringer@v0.1.12

.PHONY: generate
generate: deps
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
	go test -run XXX -bench . -benchmem ./...
