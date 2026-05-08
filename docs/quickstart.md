---
title: Quickstart
description: "From a clean machine to a 1920×1080 spectrogram in under a minute."
---

# Quickstart

One install, one command, one image.

## 1. Install

```bash
brew install steipete/tap/songsee
songsee --version
```

Other paths (go install, source builds, ffmpeg) live on [Install](install.md).

## 2. Render a default spectrogram

```bash
songsee track.mp3
```

Output is a 1920×1080 JPEG (quality 95) written next to the input — `track.jpg`. The file path is echoed to stdout; everything else (decode info under `--verbose`, warnings, errors) goes to stderr so pipes stay clean.

## 3. Pick a different view

```bash
# Mel-scaled spectrogram (perceptual frequency).
songsee track.mp3 --viz mel

# Chromagram (12-bin pitch class) with the magma palette.
songsee track.mp3 --viz chroma --style magma

# Combine harmonic/percussive split with chroma in a single grid.
songsee track.mp3 --viz hpss,chroma --style inferno
```

The full list of views lives on [Visualizations](visualizations.md); palettes are documented on [Palettes](palettes.md).

## 4. Slice a section

```bash
songsee track.mp3 --start 12.5 --duration 8 -o slice.jpg
```

`--start` and `--duration` are seconds. `-o` overrides the default output path; pass `-o -` to write the encoded image to stdout.

## 5. Stream from stdin

```bash
cat track.mp3 | songsee - --format png -o - > spectro.png
```

`-` as the input reads from stdin; `--format png` switches the encoder; `-o -` writes the encoded image to stdout. Combine with `find`, `xargs`, or shell loops for batch rendering.

## 6. Tune dimensions and FFT

```bash
songsee track.mp3 \
  --width 2560 --height 1440 \
  --window 4096 --hop 1024 \
  --min-freq 50 --max-freq 8000
```

`--window` must be a power of two. Larger windows trade time resolution for frequency resolution. `--min-freq` / `--max-freq` clamp the visible frequency band; the default upper bound is the Nyquist frequency.

## Where next

- [Visualizations](visualizations.md) — what each of the nine modes shows.
- [Palettes](palettes.md) — gradient stops for the six built-in palettes.
- [Pipeline](spec.md) — windowing, FFT, bin mapping, normalization.
- [CLI](cli.md) — every flag with its default.
