package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

// SONGSEE_BENCH_WAV selects a three-minute, 44.1 kHz stereo WAV fixture.
func BenchmarkRender(b *testing.B) {
	input := os.Getenv("SONGSEE_BENCH_WAV")
	if input == "" {
		b.Skip("set SONGSEE_BENCH_WAV to a three-minute stereo WAV")
	}
	for _, tc := range []struct {
		name string
		viz  string
	}{
		{name: "Default"},
		{name: "HPSS", viz: "hpss"},
		{name: "AllPanels", viz: "spectrogram,mel,chroma,hpss,selfsim,loudness,tempogram,mfcc,flux"},
	} {
		b.Run(tc.name, func(b *testing.B) {
			args := []string{input, "--output", "-"}
			if tc.viz != "" {
				args = append(args, "--viz", tc.viz)
			}
			var stderr bytes.Buffer
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if code := run(args, nil, io.Discard, &stderr); code != 0 {
					b.Fatalf("render exited %d: %s", code, &stderr)
				}
			}
		})
	}
}
