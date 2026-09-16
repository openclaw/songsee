package audio

import (
	"math"
	"slices"
	"testing"
)

func TestSliceRejectsNonFiniteTimes(t *testing.T) {
	a := Audio{SampleRate: 10, Samples: []float64{0, 1, 2, 3}}
	for _, value := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		for _, times := range [][2]float64{{value, 0}, {0, value}} {
			if _, err := Slice(a, times[0], times[1]); err == nil {
				t.Errorf("Slice(%v, %v) accepted non-finite time", times[0], times[1])
			}
		}
	}
}

func TestSliceLargeStart(t *testing.T) {
	a := Audio{SampleRate: 10, Samples: []float64{0, 1, 2, 3}}
	for _, start := range []float64{1e20, math.MaxFloat64} {
		if _, err := Slice(a, start, 0); err == nil {
			t.Errorf("Slice(%v, 0) accepted start beyond end", start)
		}
	}
}

func TestSliceClampsLargeDuration(t *testing.T) {
	a := Audio{SampleRate: 10, Samples: []float64{0, 1, 2, 3}}
	for _, duration := range []float64{1, 1e20, math.MaxFloat64} {
		out, err := Slice(a, 0.1, duration)
		if err != nil {
			t.Errorf("Slice(0.1, %v): %v", duration, err)
			continue
		}
		if out.SampleRate != a.SampleRate || !slices.Equal(out.Samples, a.Samples[1:]) {
			t.Errorf("Slice(0.1, %v) = %+v, want remaining samples", duration, out)
		}
	}
}
