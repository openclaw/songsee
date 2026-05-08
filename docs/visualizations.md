---
title: Visualizations
description: "The nine visualization modes songsee can render: spectrogram, mel, chroma, hpss, selfsim, loudness, tempogram, mfcc, flux."
---

# Visualizations

`--viz` selects one or more visualization modes. Pass it once with a comma-separated list, or repeat it. Unknown names error out before any decoding runs.

```bash
songsee track.mp3 --viz spectrogram
songsee track.mp3 --viz mel,chroma,hpss
songsee track.mp3 --viz spectrogram --viz flux
```

When more than one mode is selected, songsee composes a square-ish grid (`ceil(sqrt(n))` columns) with an 8 px gap between cells, all sized to fit `--width` × `--height`.

## spectrogram

Time × frequency magnitude. The base FFT view: a Hann-windowed STFT converted to decibels (`20·log10(mag + 1e-9)`). The X axis is time, the Y axis is linear frequency from `--min-freq` to `--max-freq` (Nyquist by default), and each pixel's brightness is the magnitude in that time-frequency cell.

Use it when you want raw spectral truth — verifying decode, hunting harmonics, identifying transients.

```bash
songsee track.mp3 --viz spectrogram
```

## mel

Perceptual frequency scale. Same STFT, but bins are warped onto the mel scale, which weights low frequencies more heavily — closer to how humans hear pitch.

Good for vocal and tonal content; the structure of speech and melody jumps out compared to a linear spectrogram.

```bash
songsee track.mp3 --viz mel --min-freq 80 --max-freq 8000
```

## chroma

12-bin pitch class. Energy is folded across octaves into the twelve semitones (C, C♯, D, …). The Y axis is pitch class, the X axis is time.

Reveals harmonic and key content — chord progressions, modulations, repetition between sections.

```bash
songsee track.mp3 --viz chroma
```

## hpss

Harmonic vs percussive separation. Median-filters the spectrogram twice (9-frame kernels) to split it into a harmonic top half (sustained tones) and a percussive bottom half (transients).

Use it to see where the kit and where the melody live in the same track.

```bash
songsee track.mp3 --viz hpss
```

## selfsim

Self-similarity matrix on chroma frames. Each pixel `(i, j)` is the cosine similarity between chroma frame `i` and frame `j`, with a gentle gamma (1.4) for contrast.

Brings out song structure: verses repeat as bright off-diagonal stripes; choruses form clear blocks.

```bash
songsee track.mp3 --viz selfsim
```

## loudness

Frame-wise RMS over time. A waveform-style envelope with the X axis as time and the height as energy, clamped to the 95th percentile so peaks don't crush the rest.

Good for spotting dynamics, fade-ins, and silence.

```bash
songsee track.mp3 --viz loudness
```

## tempogram

Tempo variation over time. An autocorrelation-style heatmap of the onset envelope, scanning 30–240 BPM in 256 bins.

Reveals tempo drift, rubato, and switches between rhythmic feels.

```bash
songsee track.mp3 --viz tempogram
```

## mfcc

Mel-frequency cepstral coefficients — the classic timbre fingerprint. Each row is one cepstral coefficient over time.

Strips pitch and leaves "color"; useful for distinguishing instruments, voices, or sections that share notes but not tone.

```bash
songsee track.mp3 --viz mfcc
```

## flux

Spectral flux — frame-to-frame magnitude change. A 1-D envelope with peaks at onsets and discontinuities, clamped to the 95th percentile.

Use it to find note onsets, edits, or anything sudden.

```bash
songsee track.mp3 --viz flux
```

## Combining

```bash
songsee track.mp3 --viz spectrogram,mel,chroma,hpss,selfsim,loudness,tempogram,mfcc,flux
```

All nine in one 1920×1080 grid (3×3). Mix and match with the same syntax. Each panel auto-contrasts independently — comparing absolute values across panels is not meaningful, comparing structure is.

## Related pages

- [Palettes](palettes.md) — color maps applied to every panel.
- [Pipeline](spec.md) — windowing, FFT, bin mapping, normalization.
- [CLI](cli.md) — every flag with its default.
