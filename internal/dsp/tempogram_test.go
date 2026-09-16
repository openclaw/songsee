package dsp

import "testing"

func TestTempogramPreservesTempoWhenCapped(t *testing.T) {
	// A 120 BPM onset every 50 frames at 100 frames/second.
	spec := Spectrogram{Frames: 1200, Bins: 1, SampleRate: 8000, HopSize: 80, Values: make([]float64, 1200)}
	for frame := 25; frame < spec.Frames; frame += 50 {
		spec.Values[frame] = 1
	}
	full := Tempogram(&spec, 100, 140, 0)
	if full.At(600, 20) == 0 || full.At(600, 20) <= full.At(600, 0) {
		t.Fatal("fixture must have stronger 120 BPM than 100 BPM correlation")
	}
	for _, limit := range []int{1, 128, 256, 1200, 1300} {
		capped := Tempogram(&spec, 100, 140, limit)
		wantWidth := min(limit, spec.Frames)
		if capped.Width != wantWidth {
			t.Fatalf("limit %d: width %d, want %d", limit, capped.Width, wantWidth)
		}
		for x := 0; x < capped.Width; x++ {
			center := spec.Frames / 2
			if capped.Width > 1 {
				center = int(float64(x) * float64(spec.Frames-1) / float64(capped.Width-1))
			}
			for y := 0; y < capped.Height; y++ {
				if got, want := capped.At(x, y), full.At(center, y); got != want {
					t.Fatalf("limit %d, column %d, BPM %d: correlation %v, want %v at original frame %d", limit, x, y+100, got, want, center)
				}
			}
		}
	}
}

func TestTempogramEmpty(t *testing.T) {
	spec := Spectrogram{SampleRate: 8000, HopSize: 80}
	out := Tempogram(&spec, 100, 140, 256)
	if out.Width != 0 || len(out.Values) != 0 {
		t.Fatal("expected an empty tempogram")
	}
}
