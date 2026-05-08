---
title: Palettes
description: "Six built-in color maps for songsee: classic, magma, inferno, viridis, gray, clawd."
---

# Palettes

`--style` picks a palette. All palettes are 5- or 6-stop linear gradients applied to normalized values in `[0, 1]`. The default is `classic`.

```bash
songsee track.mp3 --style magma
songsee track.mp3 --viz mel --style viridis
songsee track.mp3 --viz hpss,chroma --style clawd
```

Unknown names error out before decoding. All palettes are deterministic — the same input always produces the same colors.

## classic

The default. A black → navy → cyan → amber → white sweep tuned for log-magnitude data, with strong perceptual contrast across the full range.

| stop | color |
|------|-------|
| 0.00 | `#000000` |
| 0.20 | `#002060` |
| 0.45 | `#00a0c8` |
| 0.70 | `#ffb400` |
| 1.00 | `#ffffff` |

## magma

Matplotlib's magma. Black → deep purple → magenta → orange → cream. Smooth and perceptually uniform; works well for everything from spectrograms to MFCCs.

| stop | color |
|------|-------|
| 0.00 | `#000004` |
| 0.25 | `#3b0c57` |
| 0.50 | `#b4367a` |
| 0.75 | `#fb8c3c` |
| 1.00 | `#fcfdbf` |

## inferno

Matplotlib's inferno. Same shape as magma with hotter highs — black → indigo → red → orange → pale yellow.

| stop | color |
|------|-------|
| 0.00 | `#000004` |
| 0.25 | `#3d0965` |
| 0.50 | `#bb3754` |
| 0.75 | `#f98e08` |
| 1.00 | `#fcffa4` |

## viridis

Matplotlib's viridis. Purple → blue → teal → green → yellow. Colorblind-safe and perceptually uniform; the safest choice for publication figures.

| stop | color |
|------|-------|
| 0.00 | `#440154` |
| 0.25 | `#3a528b` |
| 0.50 | `#20908c` |
| 0.75 | `#5ec962` |
| 1.00 | `#fde725` |

## gray

A straight black-to-white linear ramp. Ideal for print, monochrome compositing, or downstream processing that doesn't want hue information.

| stop | color |
|------|-------|
| 0.00 | `#000000` |
| 1.00 | `#ffffff` |

`grey` is accepted as an alias.

## clawd 🦞

The mascot palette. Abyssal navy → ocean teal → coral → lobster red → foam highlight. Six stops, designed to be unmistakable.

| stop | color |
|------|-------|
| 0.00 | `#02040f` |
| 0.20 | `#0b264a` |
| 0.40 | `#126175` |
| 0.60 | `#c1625c` |
| 0.80 | `#cd3728` |
| 1.00 | `#ffe6d2` |

## Related pages

- [Visualizations](visualizations.md) — the nine viz modes that consume these palettes.
- [Pipeline](spec.md) — how values are normalized before palette mapping.
