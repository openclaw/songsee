VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS ?= -s -w -X main.version=$(VERSION)

.PHONY: songsee docs-site

songsee:
	mkdir -p bin
	go build -ldflags "$(LDFLAGS)" -o bin/songsee ./cmd/songsee

docs-site:
	@node scripts/build-docs-site.mjs
