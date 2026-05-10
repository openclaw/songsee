// Package render turns spectrograms into images.
package render

import (
	"fmt"
	"image"

	"github.com/steipete/songsee/internal/dsp"
)

// HeatmapOptions configures generic feature map rendering.
type HeatmapOptions struct {
	Width    int
	Height   int
	Palette  Palette
	Min      float64
	Max      float64
	Clamp    bool
	FlipVert bool
}

// Heatmap renders a feature map into an RGBA image.
func Heatmap(mapIn *dsp.FeatureMap, opts HeatmapOptions) (*image.RGBA, error) {
	if mapIn == nil {
		return nil, fmt.Errorf("feature map required")
	}
	if opts.Width <= 0 || opts.Height <= 0 {
		return nil, fmt.Errorf("invalid output size")
	}
	if opts.Palette == nil {
		return nil, fmt.Errorf("palette required")
	}
	if mapIn.Width <= 0 || mapIn.Height <= 0 {
		return nil, fmt.Errorf("invalid feature map")
	}

	minVal := mapIn.Min
	maxVal := mapIn.Max
	if opts.Clamp {
		minVal = opts.Min
		maxVal = opts.Max
	}
	if maxVal <= minVal {
		maxVal = minVal + 1
	}

	img := image.NewRGBA(image.Rect(0, 0, opts.Width, opts.Height))
	for x := 0; x < opts.Width; x++ {
		srcX := sampleIndex(x, opts.Width, mapIn.Width)
		for y := 0; y < opts.Height; y++ {
			srcY := sampleIndex(y, opts.Height, mapIn.Height)
			if opts.FlipVert {
				srcY = mapIn.Height - 1 - srcY
			}
			val := mapIn.Values[srcY*mapIn.Width+srcX]
			norm := (val - minVal) / (maxVal - minVal)
			if norm < 0 {
				norm = 0
			}
			if norm > 1 {
				norm = 1
			}
			setRGBA(img, x, y, opts.Palette(norm))
		}
	}
	return img, nil
}
