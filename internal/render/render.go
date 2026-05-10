// Package render turns spectrograms into images.
package render

import (
	"fmt"
	"image"
	"math"

	"github.com/steipete/songsee/internal/dsp"
)

// Options configures spectrogram rendering.
type Options struct {
	Width    int
	Height   int
	MinFreq  float64
	MaxFreq  float64
	Palette  Palette
	MinDB    float64
	MaxDB    float64
	ClampDB  bool
	FlipVert bool
}

// Spectrogram renders a spectrogram into an RGBA image.
func Spectrogram(spec *dsp.Spectrogram, opts Options) (*image.RGBA, error) {
	if spec == nil {
		return nil, fmt.Errorf("spectrogram required")
	}
	if opts.Width <= 0 || opts.Height <= 0 {
		return nil, fmt.Errorf("invalid output size")
	}
	if opts.Palette == nil {
		return nil, fmt.Errorf("palette required")
	}

	minDB := spec.Min
	maxDB := spec.Max
	if opts.ClampDB {
		minDB = opts.MinDB
		maxDB = opts.MaxDB
	}
	if maxDB <= minDB {
		maxDB = minDB + 1
	}

	minBin := 0
	maxBin := spec.Bins - 1
	if opts.MinFreq > 0 {
		minBin = int(opts.MinFreq / spec.BinHz)
	}
	if opts.MaxFreq > 0 {
		maxBin = int(opts.MaxFreq / spec.BinHz)
	}
	if minBin < 0 {
		minBin = 0
	}
	if maxBin >= spec.Bins {
		maxBin = spec.Bins - 1
	}
	if maxBin <= minBin {
		minBin = 0
		maxBin = spec.Bins - 1
	}
	binSpan := maxBin - minBin

	img := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
	frames := spec.Frames
	bins := spec.Bins
	for x := 0; x < opts.Width; x++ {
		frame := sampleIndex(x, opts.Width, frames)
		frameOffset := frame * bins
		for y := 0; y < opts.Height; y++ {
			pos := 0.0
			if opts.Height > 1 {
				pos = float64(y) / float64(opts.Height-1)
			}
			bin := minBin + int(math.Round((1-pos)*float64(binSpan)))
			if bin < minBin {
				bin = minBin
			}
			if bin > maxBin {
				bin = maxBin
			}
			val := spec.Values[frameOffset+bin]
			norm := (val - minDB) / (maxDB - minDB)
			if norm < 0 {
				norm = 0
			}
			if norm > 1 {
				norm = 1
			}
			c := opts.Palette(norm)
			ypos := y
			if opts.FlipVert {
				ypos = opts.Height - 1 - y
			}
			setRGBA(img, x, ypos, c)
		}
	}
	return img, nil
}

func sampleIndex(pos, dstSize, srcSize int) int {
	if srcSize <= 1 || dstSize <= 1 {
		return 0
	}
	return int(math.Round(float64(pos) * float64(srcSize-1) / float64(dstSize-1)))
}
