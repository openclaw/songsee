package dsp

import (
	"math"
	"testing"
)

func BenchmarkSelfSimilarity(b *testing.B) {
	// The CLI caps self-similarity at 200 frames of 12 chroma features.
	features := NewFeatureMap(200, 12)
	for x := 0; x < features.Width; x++ {
		for y := 0; y < features.Height; y++ {
			features.Set(x, y, 20*math.Sin(float64(x)/17+float64(y)))
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SelfSimilarity(features, 200)
	}
}

func BenchmarkHPSS(b *testing.B) {
	// A three-second excerpt with the CLI's 2048-sample window and 512-sample hop.
	samples := make([]float64, 3*44100)
	for i := range samples {
		t := float64(i) / 44100
		samples[i] = 0.5*math.Sin(2*math.Pi*220*t) +
			0.2*math.Sin(2*math.Pi*330*t) +
			0.3*math.Exp(-40*math.Mod(t, 0.5))*math.Sin(2*math.Pi*80*t)
	}
	spec := ComputeSpectrogram(samples, 44100, 2048, 512)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HPSS(&spec, 9, 9)
	}
}
