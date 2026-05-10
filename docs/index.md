---
title: Overview
permalink: /
description: "songsee is a single Go CLI that turns audio into modern spectrogram and feature-panel images — fast WAV/MP3 decode, ffmpeg fallback, nine visualization modes, six palettes."
---

## Try it

After [installing](install.md), every render is a one-liner.

```bash
# Default: a clean spectrogram next to the input file.
songsee track.mp3

# Mel spectrogram, magma palette, 2K wide.
songsee track.mp3 --viz mel --style magma --width 2048 --height 1024

# All nine modes in one grid.
songsee track.mp3 --viz spectrogram,mel,chroma,hpss,selfsim,loudness,tempogram,mfcc,flux

# Slice eight seconds out of a long file.
songsee track.mp3 --start 12.5 --duration 8 -o slice.jpg

# Pipe from stdin, write PNG to stdout.
cat track.mp3 | songsee - --format png -o - > spectro.png
```

The default output is a 1920×1080 JPEG (quality 95) written next to the input. `--format png` switches encoder, `-o` redirects the path, and `-o -` streams to stdout for piping.

## What songsee does

- **One binary, nine views.** spectrogram, mel, chroma, hpss, selfsim, loudness, tempogram, mfcc, flux — pick one, combine several, or render the full grid.
- **Fast decode paths.** Native Go decoders for WAV (PCM, float, extensible) and MP3; ffmpeg fallback covers everything else.
- **Six palettes.** classic, magma, inferno, viridis, gray, and claw — each tuned for log-magnitude data.
- **Auto-contrast.** Per-panel percentile clamping (0.05 / 0.98) keeps every visualization readable without manual tuning.
- **Scriptable I/O.** File path, stdin (`-`), or stdout. Quiet mode for CI; verbose mode prints decode and slice details to stderr.
- **No Python.** Single static binary. No model files, no virtualenv, no GPU.

## Pick your path

- **Trying it.** [Install](install.md) → [Quickstart](quickstart.md). One brew formula, one command, one image.
- **Picking a view.** [Visualizations](visualizations.md) describes what each of the nine modes shows and when to use it.
- **Picking a palette.** [Palettes](palettes.md) lists the six palettes with their gradient stops.
- **Audio inputs.** [Decoding](decoding.md) covers WAV/MP3 fast paths, ffmpeg fallback, sample rate, and stdin.
- **Output and batches.** [Rendering](rendering.md) explains output sizing, grid layout, format selection, and stdout streaming.
- **Algorithm details.** [Pipeline](spec.md) documents windowing, FFT, bin mapping, and normalization.
- **Flag reference.** [CLI](cli.md) lists every flag with its default.

## Project

Active development; the [changelog](https://github.com/openclaw/songsee/blob/main/CHANGELOG.md) tracks what shipped. Released under the [MIT license](https://github.com/openclaw/songsee/blob/main/LICENSE). Source on [GitHub](https://github.com/openclaw/songsee).
