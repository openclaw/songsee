---
title: CLI
description: "Every songsee flag with its default and accepted values."
---

# CLI

```text
songsee <input> [flags]
```

`<input>` is a file path or `-` for stdin.

## Inputs and output

| flag | type | default | description |
|------|------|---------|-------------|
| `<input>` | string (positional) | required | File path, or `-` to read encoded audio from stdin. |
| `-o`, `--output` | string | input name + extension | Output path. `-` writes the encoded image to stdout. |
| `--format` | `jpg` \| `png` | `jpg` | Output encoder. JPEG is quality 95; PNG is lossless. |
| `--width` | int | `1920` | Output width in pixels. |
| `--height` | int | `1080` | Output height in pixels. |
| `-q`, `--quiet` | bool | `false` | Suppress the stdout output-path echo. |
| `-v`, `--verbose` | bool | `false` | Print decode and slice info to stderr. |
| `--version` | bool | — | Print version and exit. |

## FFT and windowing

| flag | type | default | description |
|------|------|---------|-------------|
| `--window` | int | `2048` | FFT window size in samples. Must be a power of two. |
| `--hop` | int | `512` | Hop size in samples between frames. |
| `--min-freq` | float (Hz) | `0` | Lower bound of the visible frequency band. |
| `--max-freq` | float (Hz) | Nyquist | Upper bound of the visible frequency band. Must exceed `--min-freq`. |

## Slicing

| flag | type | default | description |
|------|------|---------|-------------|
| `--start` | float (s) | `0` | Skip this many seconds from the start of the input. |
| `--duration` | float (s) | `0` (full) | Render only this many seconds after `--start`. |

## Visualization

| flag | type | default | description |
|------|------|---------|-------------|
| `--viz` | repeated string list | `spectrogram` | One or more of: `spectrogram`, `mel`, `chroma`, `hpss`, `selfsim`, `loudness`, `tempogram`, `mfcc`, `flux`. Repeatable or comma-separated. |
| `--style` | string | `classic` | Palette name: `classic`, `magma`, `inferno`, `viridis`, `gray` (alias `grey`), `claw` (legacy alias `clawd`). |

## Decoding

| flag | type | default | description |
|------|------|---------|-------------|
| `--sample-rate` | int | `44100` | Sample rate requested from the ffmpeg fallback. Native WAV/MP3 keep the file's rate. |
| `--ffmpeg` | string | first `ffmpeg` on `PATH` | Override the ffmpeg binary used for non-WAV/MP3 inputs. |

## Exit codes

| code | meaning |
|------|---------|
| `0` | Render succeeded. |
| `1` | Decode, render, or write error (message on stderr). |
| `2` | Usage error — bad flag, invalid combination, or missing input. |

## Examples

```bash
# All defaults.
songsee track.mp3

# Mel + chroma in viridis at 2K.
songsee track.mp3 --viz mel,chroma --style viridis --width 2048 --height 1024

# Eight-second slice starting at 12.5s, written to PNG.
songsee track.mp3 --start 12.5 --duration 8 -o slice.png

# Stream from stdin, encode to PNG, write to stdout.
cat track.mp3 | songsee - --format png -o - > spectro.png

# Custom FFT, sub-bass focus.
songsee track.mp3 --window 4096 --hop 1024 --min-freq 20 --max-freq 200

# Pin a specific ffmpeg.
songsee weird.opus --ffmpeg /opt/homebrew/bin/ffmpeg
```

## Related pages

- [Quickstart](quickstart.md) — first render in under a minute.
- [Visualizations](visualizations.md) — when to use each viz mode.
- [Palettes](palettes.md) — palette gradient stops.
- [Decoding](decoding.md) — input formats and ffmpeg fallback.
- [Rendering](rendering.md) — output sizing, format, batch use.
