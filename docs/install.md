---
title: Install
description: "Install songsee via Homebrew, go install, or build from source."
---

# Install

`songsee` ships as a single Go binary. Pick whichever delivery mechanism fits.

## Homebrew (macOS, Linux)

```bash
brew install steipete/tap/songsee
songsee --version
```

The formula lives in [`steipete/homebrew-tap`](https://github.com/steipete/homebrew-tap). `brew upgrade songsee` brings in new releases.

## go install

```bash
go install github.com/steipete/songsee/cmd/songsee@latest
songsee --version
```

This builds against the Go version declared in `go.mod`. The binary lands in `$(go env GOBIN)` (or `$(go env GOPATH)/bin`).

## Build from source

```bash
git clone https://github.com/openclaw/songsee.git
cd songsee
make
./songsee --version
```

`make` runs `go build` with the version string injected from `git describe`.

## ffmpeg (optional)

WAV and MP3 decode natively in pure Go. Anything else (FLAC, AAC, OGG, M4A, video containers) falls through to `ffmpeg` on `PATH`.

```bash
brew install ffmpeg          # macOS / Linuxbrew
apt install ffmpeg           # Debian / Ubuntu
```

Override the lookup with `--ffmpeg /custom/path/ffmpeg` when you have several builds installed.

## Verify the install

```bash
songsee --version
songsee --help
songsee testdata/short.wav   # render a tiny known-good file
```

## Updating

- **Homebrew:** `brew upgrade songsee`.
- **go install:** rerun `go install github.com/steipete/songsee/cmd/songsee@latest`.
- **Source:** `git pull && make` — version comes from `git describe`.

## Related pages

- [Quickstart](quickstart.md) — first render in under a minute.
- [Decoding](decoding.md) — WAV/MP3 fast paths and ffmpeg fallback.
- [CLI](cli.md) — every flag with its default.
