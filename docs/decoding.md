---
title: Decoding
description: "How songsee decodes audio — native WAV/MP3 paths, ffmpeg fallback, sample rate, stdin, and slicing."
---

# Decoding

songsee turns the input into mono `float64` samples before any analysis runs. Two fast paths cover most files; everything else falls through to ffmpeg.

## Inputs

- **File path.** `songsee track.mp3` — any path the OS can open.
- **Stdin.** `songsee -` — reads the encoded stream from stdin. Useful behind `cat`, `curl`, or shell pipelines.
- **Mono mixdown.** Stereo or multichannel inputs are averaged to mono before windowing.

## Native WAV

Pure-Go WAV decoder. Handles:

- PCM 8/16/24/32-bit integer
- 32-bit float, 64-bit float
- WAVE_FORMAT_EXTENSIBLE (with channel masks and sub-format GUIDs)

No external dependency, no ffmpeg roundtrip. The decoder validates the RIFF header and rejects truncated `data` chunks before allocating sample buffers.

## Native MP3

Pure-Go MP3 decoder. Handles MPEG-1/2 Layer III with VBR and CBR. Output sample rate is whatever the file declares; songsee does not resample.

If the decoder hits a malformed frame it surfaces a structured error instead of silently truncating, so corrupt input fails loudly.

## ffmpeg fallback

Anything that isn't WAV or MP3 — FLAC, AAC, M4A, OGG, Opus, video containers, raw streams — is decoded by spawning `ffmpeg`. songsee asks for 32-bit float little-endian mono at the configured sample rate:

```text
ffmpeg -hide_banner -loglevel error \
       -i <input> -f f32le -ac 1 -ar <sample-rate> -
```

Tweak the pipeline with:

- `--sample-rate N` — output sample rate fed to ffmpeg (default `44100`).
- `--ffmpeg /path/to/ffmpeg` — override the binary lookup.

If ffmpeg isn't on `PATH` and the file isn't WAV or MP3, songsee fails with a clear error. Install with `brew install ffmpeg` or your distro's package manager.

## Slicing

`--start` and `--duration` slice the decoded audio before analysis. Both are seconds (float). Negative values are rejected.

```bash
songsee long.mp3 --start 60 --duration 15 -o minute1.jpg
songsee long.mp3 --start 60                  # 60s to end
songsee long.mp3 --duration 30               # first 30s
```

Slicing happens on samples after decoding; FFT framing then runs on the slice.

## Sample rate notes

- Native WAV/MP3 keep the file's sample rate. `--sample-rate` only affects the ffmpeg pipeline.
- The Nyquist frequency (`sampleRate / 2`) is the upper bound visible in the spectrogram.
- 44.1 kHz / 48 kHz inputs render with the default `--window 2048 --hop 512` at ≈21 ms / frame.

## Verbose decoding

```bash
songsee track.flac --verbose -o out.png
```

`--verbose` prints decode info to stderr — sample count, sample rate, slice bounds — without polluting stdout. Combine with `--quiet` to suppress the trailing output-path echo when piping.

## Related pages

- [Install](install.md) — how to add ffmpeg if you need it.
- [Pipeline](spec.md) — what happens to the samples after decoding.
- [CLI](cli.md) — every flag with its default.
