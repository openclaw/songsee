package audio

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestNonLayerIIIAudioUsesFFmpeg(t *testing.T) {
	for _, tt := range []struct {
		name   string
		header byte
	}{
		{"MPEG4-ADTS.aac", 0xf1},
		{"MPEG2-ADTS.aac", 0xf9},
		{"MPEG1-LayerI.mp1", 0xff},
		{"MPEG1-LayerII.mp2", 0xfd},
	} {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, 16)
			copy(data, []byte{0xff, tt.header, 0x6c, 0x40})
			reader := bytes.NewReader(data)
			if _, ok, err := DecodeMP3If(reader); ok || err != nil {
				t.Errorf("non-Layer-III header recognized as MP3: ok=%v err=%v", ok, err)
			}
			if pos, err := reader.Seek(0, io.SeekCurrent); err != nil || pos != 0 {
				t.Errorf("probe did not rewind: pos=%d err=%v", pos, err)
			}
			opts := Options{SampleRate: 22050, FFmpegPath: installFakeFFmpeg(t)}
			pcm, err := DecodeBytes(data, opts)
			if err != nil || pcm.SampleRate != 22050 || !slices.Equal(pcm.Samples, []float64{0, 0.5}) {
				t.Errorf("byte input did not use ffmpeg: pcm=%+v err=%v", pcm, err)
			}
			path := filepath.Join(t.TempDir(), tt.name)
			if err := os.WriteFile(path, data, 0o644); err != nil {
				t.Fatal(err)
			}
			pcm, err = DecodeFile(path, opts)
			if err != nil || pcm.SampleRate != 22050 || !slices.Equal(pcm.Samples, []float64{0, 0.5}) {
				t.Errorf("file input did not use ffmpeg: pcm=%+v err=%v", pcm, err)
			}
		})
	}
}

func TestCorruptLayerIIIStillReturnsDecodeError(t *testing.T) {
	data := []byte{0xff, 0xfb, 0, 0, 0, 0, 0, 0}
	if _, ok, err := DecodeMP3If(bytes.NewReader(data)); !ok || err == nil {
		t.Fatalf("corrupt Layer III must remain a native decode error: ok=%v err=%v", ok, err)
	}
}
