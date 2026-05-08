.PHONY: songsee docs-site

songsee:
	mkdir -p bin
	go build -o bin/songsee ./cmd/songsee

docs-site:
	@node scripts/build-docs-site.mjs
