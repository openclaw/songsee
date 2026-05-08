---
title: Pipeline
description: "songsee's spectral pipeline — decode, window, FFT, bin mapping, normalization, render."
---

# Pipeline

This page documents the algorithm and the defaults songsee uses to produce repeatable, high-quality images. It complements [Visualizations](visualizations.md) (what each mode shows) and [Rendering](rendering.md) (how the canvas gets composed).

## Stages

```text
input → decode → mono mixdown → optional slice → window → FFT
       ↓
       per-mode features (mel, chroma, mfcc, hpss, …)
       ↓
       percentile normalize → palette map → grid compose → encode
```

Every stage is deterministic. The same input file with the same flags always produces the same output bytes.

## Decode

- WAV (PCM 8/16/24/32-bit, 32/64-bit float, WAVE_FORMAT_EXTENSIBLE) and MP3 are decoded in pure Go via the bundled decoders.
- Anything else falls through to `ffmpeg` (32-bit float little-endian, mono, `--sample-rate` Hz; default `44100`).
- Stereo or multichannel input is averaged to mono.
- `--start` / `--duration` slice the decoded sample buffer in seconds before windowing.

See [Decoding](decoding.md) for input formats, sample rate, ffmpeg lookup, and stdin usage.

## Windowing and FFT

- Window: **Hann**, applied per frame.
- Window size: `--window` samples (default `2048`, must be a power of two).
- Hop size: `--hop` samples (default `512`).
- Frame count: `1 + (len(samples) - window + hop - 1) / hop`.
- Bin count: `window / 2 + 1`.
- Bin spacing: `sampleRate / window` Hz per bin.

Magnitude is converted to decibels with `20·log10(mag + 1e-9)` for the base spectrogram. Per-feature pipelines (mel, chroma, mfcc) use linear power instead.

## Per-mode features

| mode | source | notes |
|------|--------|-------|
| `spectrogram` | STFT magnitude in dB | clamped to 5th–98th percentile |
| `mel` | mel-warped power | log-magnitude; clamped 5th–98th percentile |
| `chroma` | 12-bin pitch class | folds octaves; clamped 10th–98th percentile |
| `mfcc` | DCT of mel power | strips pitch, keeps timbre |
| `hpss` | median filters on STFT | 9-frame harmonic + 9-frame percussive kernels |
| `selfsim` | cosine sim on chroma frames | gamma 1.4; clamped 10th–98th percentile |
| `loudness` | per-frame RMS | clamped to 95th percentile |
| `tempogram` | onset autocorrelation | 30–240 BPM, 256 bins |
| `flux` | frame-to-frame STFT delta | clamped to 95th percentile |

The percentile sampling reservoir is capped at 20 000 values per panel for speed; this is dense enough that boundaries are stable across runs.

## Rendering

- Each panel maps `(time × bin)` cells onto pixels at the panel's width × height.
- Values are normalized into `[0, 1]` against the per-panel min/max (after the percentile clamp), then passed through the chosen palette.
- Heatmap panels (mel, chroma, mfcc, selfsim, hpss halves, tempogram) render with `flipVert` so low frequencies are at the bottom.
- Multiple panels compose into a `ceil(sqrt(n))`-column grid with an 8 px gap (see [Rendering](rendering.md)).
- Encoder: PNG (lossless) or JPEG (quality 95).

## CLI defaults

```text
--format       jpg
--width        1920
--height       1080
--window       2048
--hop          512
--sample-rate  44100
--style        classic
--viz          spectrogram
```

Full reference: [CLI](cli.md).

## Related pages

- [Visualizations](visualizations.md) — per-mode descriptions.
- [Palettes](palettes.md) — gradient stops.
- [Decoding](decoding.md) — input handling.
- [Rendering](rendering.md) — output and batch use.
