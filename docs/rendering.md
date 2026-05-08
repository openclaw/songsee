---
title: Rendering
description: "How songsee turns spectrogram data into images — output sizing, format, grid layout, stdout streaming, and batch use."
---

# Rendering

The render stage maps numeric spectrogram and feature data onto pixels, applies the chosen palette, composes panels into a grid, and encodes the result to JPEG or PNG.

## Output path

```bash
songsee track.mp3                       # writes track.jpg next to the input
songsee track.mp3 -o out.png            # explicit path; format inferred from extension
songsee track.mp3 -o spectro            # no extension; appends ".jpg" by default
songsee - -o -                          # stdin in, encoded image to stdout
```

If `--format` is set explicitly, it overrides extension-based inference. If `-o` already ends in `.png`, `.jpg`, or `.jpeg`, the encoder follows the extension regardless of `--format`.

When the input is `-` (stdin) and no `-o` is given, the output filename is `songsee.jpg` (or `.png`) in the current directory.

## Format

```bash
songsee track.mp3 --format png          # PNG, lossless
songsee track.mp3 --format jpg          # JPEG, quality 95 (default)
```

JPEG quality is fixed at 95 — high enough that compression artifacts disappear at typical viewing sizes while keeping file sizes reasonable. PNG is the right choice for archival, transparency, or downstream processing that doesn't tolerate JPEG quantization.

## Dimensions

```bash
songsee track.mp3 --width 2560 --height 1440
songsee track.mp3 --width 3840 --height 2160       # 4K
songsee track.mp3 --width 800 --height 200         # banner strip
```

Defaults are 1920×1080. Both must be positive. With multiple visualizations, songsee divides the canvas into a grid; very small canvases with many panels can leave cells too small to render and produce an error.

## Grid layout

When `--viz` selects more than one mode, panels are tiled into a `ceil(sqrt(n))`-column grid with an 8 px gap, sized to fit `--width` × `--height` exactly:

| panels | grid |
|--------|------|
| 1 | 1×1 |
| 2 | 2×1 |
| 3, 4 | 2×2 |
| 5–6 | 3×2 |
| 7–9 | 3×3 |

Cells are equal width and height; the canvas is filled top-down, left-right, in the order panels appear in `--viz`.

## Frequency range

`--min-freq` and `--max-freq` clamp the visible band (Hz) for spectrogram, mel, and MFCC panels. The default upper bound is the Nyquist frequency (`sampleRate / 2`).

```bash
# Vocal range, 80 Hz – 4 kHz.
songsee track.mp3 --min-freq 80 --max-freq 4000

# Sub-bass focus.
songsee track.mp3 --min-freq 20 --max-freq 200
```

`--max-freq` must be greater than `--min-freq`; songsee rejects the run with exit code 2 otherwise.

## Auto-contrast

Every panel runs an independent percentile clamp on its values before palette mapping (typically 0.05 / 0.98 — ~5% black floor, ~2% white ceiling). This keeps a quiet ambient track and a loud rock track equally readable in the same grid; it also means absolute brightness is not comparable across panels.

The base spectrogram converts magnitudes to decibels (`20·log10(mag + 1e-9)`) before normalizing.

## Stdout streaming

Pass `-o -` to write the encoded image bytes to stdout. Combine with `--quiet` to silence the trailing path echo:

```bash
songsee track.mp3 -o - --quiet > spectro.jpg
songsee track.mp3 -o - --format png --quiet | imgcat
ssh host "songsee /audio/x.flac -o -" > x.jpg
```

Verbose decode info still goes to stderr under `--verbose`, so it doesn't corrupt the binary stream on stdout.

## Batch usage

There's no built-in `songsee batch`; lean on the shell.

```bash
# All MP3s in a directory.
for f in *.mp3; do songsee "$f" --style magma; done

# Parallel via xargs.
ls *.mp3 | xargs -P 8 -I{} songsee {} --width 1920 --height 540

# find + GNU parallel.
find . -name '*.flac' -print0 | parallel -0 songsee {} --style viridis -o {.}.png
```

songsee is single-threaded internally, so parallelism comes from running multiple processes.

## Related pages

- [Visualizations](visualizations.md) — what each panel shows.
- [Palettes](palettes.md) — color maps applied at render time.
- [Pipeline](spec.md) — windowing, FFT, normalization details.
- [CLI](cli.md) — every flag with its default.
