BINARY    := codeowner
MODULE    := github.com/kevin-robayna/codeowner
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT    := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE      := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS   := -s -w \
	-X $(MODULE)/internal/appinfo.Version=$(VERSION) \
	-X $(MODULE)/internal/appinfo.Commit=$(COMMIT) \
	-X $(MODULE)/internal/appinfo.Date=$(DATE)

.PHONY: build test lint lint-check clean

build: clean
	go build -trimpath -ldflags '$(LDFLAGS)' -o $(BINARY) ./cmd/codeowner

test:
	go test -race -cover ./...

lint:
	golangci-lint run --fix

lint-check:
	golangci-lint run

clean:
	rm -f $(BINARY)
	rm -rf dist/
