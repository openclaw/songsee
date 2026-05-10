# Changelog

## 0.1.1 - 2026-05-10

- New Clawd style
- Docs: rewritten gogcli-style — plain-markdown pages for install, quickstart, visualizations, palettes, decoding, rendering, pipeline, and CLI; new custom static-site builder (`make docs-site`) and `pages.yml` workflow render songsee.sh with a sidebar nav, search, dark-mode toggle, and per-page TOC
- Source builds: `make` now injects the `git describe` version string
- Dependencies: updated Kong to 1.15.0

## 0.1.0 - 2026-01-02

- Spectrogram + feature panels (mel, chroma, hpss, selfsim, loudness, tempogram, mfcc, flux) with multi-panel grid
- Native WAV/MP3 decoding (ffmpeg fallback for other formats)
- PNG/JPEG output with size control and time slicing
- Palette styles: classic, magma, inferno, viridis, gray
