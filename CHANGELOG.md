# Changelog

## 0.1.3 - Unreleased

## 0.1.2 - 2026-09-05

**Highlights:** Run Songsee in Docker with ffmpeg included, and get accurate version reporting from Go-installed builds.

- Docker: add a local image with bundled `ffmpeg` for batch and server rendering without host audio tooling.
- Decoding: reject truncated WAV format/data payloads before allocation or seeking, preserve complete large WAV files, and add an opt-in `--ffmpeg-timeout` with an unlimited default. Thanks @SebTardif.
- Install: make `go install ...@latest` builds report their tagged module version instead of `dev`.
- CLI: list the `claw` palette in `--help` and reject invalid `--style`/`--viz` values before decoding input.
- Palettes: rename `clawd` to `claw`; the old name remains accepted as a compatibility alias.
- Rendering: reduce per-pixel overhead in heatmap, spectrogram, and curve drawing.
- Docs: fix table-of-contents links for formatted and repeated headings, correct the install smoke-test fixture, and align the release checklist with automated Homebrew updates. Thanks @vincentkoc.
- Dependencies: update Kong to 1.16.1.
- Build and CI: update to Go 1.27, Node 26, Alpine 3.24, Dockerfile frontend 1.27, current analysis tools, and pinned GitHub Actions; add deadcode checks, docs tests, Docker smoke coverage, and restricted workflow token permissions. Thanks @vincentkoc.

## 0.1.1 - 2026-05-10

- New Claw style
- Docs: rewritten gogcli-style — plain-markdown pages for install, quickstart, visualizations, palettes, decoding, rendering, pipeline, and CLI; new custom static-site builder (`make docs-site`) and `pages.yml` workflow render songsee.sh with a sidebar nav, search, dark-mode toggle, and per-page TOC
- Source builds: `make` now injects the `git describe` version string
- Dependencies: updated Kong to 1.15.0

## 0.1.0 - 2026-01-02

- Spectrogram + feature panels (mel, chroma, hpss, selfsim, loudness, tempogram, mfcc, flux) with multi-panel grid
- Native WAV/MP3 decoding (ffmpeg fallback for other formats)
- PNG/JPEG output with size control and time slicing
- Palette styles: classic, magma, inferno, viridis, gray
