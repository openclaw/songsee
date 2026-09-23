package audio

import (
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
)

func TestDecodeWAVRejectsPartialFrames(t *testing.T) {
	formats := []struct {
		format uint16
		bits   int
	}{
		{1, 8}, {1, 16}, {1, 24}, {1, 32}, {3, 32}, {3, 64},
	}
	for _, format := range formats {
		for _, channels := range []int{1, 2, 3} {
			frameSize := channels * format.bits / 8
			if frameSize == 1 {
				continue
			}
			t.Run(fmt.Sprintf("format%d/%dbit/%dchannels", format.format, format.bits, channels), func(t *testing.T) {
				data := makeWAVCustom(format.format, format.bits, make([]byte, 2*frameSize-1), 44100)
				binary.LittleEndian.PutUint16(data[22:24], uint16(channels))
				binary.LittleEndian.PutUint32(data[28:32], uint32(44100*frameSize))
				binary.LittleEndian.PutUint16(data[32:34], uint16(frameSize))
				// RIFF padding is outside the declared audio payload.
				data = append(data, 0)
				binary.LittleEndian.PutUint32(data[4:8], uint32(len(data)-8))
				_, err := DecodeBytes(data, Options{})
				if err == nil || !strings.Contains(err.Error(), "incomplete sample frame") {
					t.Fatalf("expected incomplete sample frame error, got %v", err)
				}
			})
		}
	}
}

func TestDecodeWAVRejectsZeroSampleRate(t *testing.T) {
	_, err := DecodeBytes(makeWAV([]int16{1, 2}, 0, 1), Options{})
	if err == nil || !strings.Contains(err.Error(), "invalid sample rate") {
		t.Fatalf("expected invalid sample rate error, got %v", err)
	}
}

func TestDecodeWAVRejectsUnsupportedEmptyFormats(t *testing.T) {
	for _, format := range []struct {
		format uint16
		bits   int
	}{
		{1, 12}, {1, 64}, {3, 16}, {3, 24},
	} {
		t.Run(fmt.Sprintf("format%d/%dbit", format.format, format.bits), func(t *testing.T) {
			_, err := DecodeBytes(makeWAVCustom(format.format, format.bits, nil, 44100), Options{})
			if err == nil || !strings.Contains(err.Error(), "unsupported bit depth") {
				t.Fatalf("expected unsupported bit depth error, got %v", err)
			}
		})
	}
}

func TestDecodeWAVKeepsOddCompletePayload(t *testing.T) {
	data := makeWAVPCM(24, []int32{1, 2, 3}, 44100)
	data = append(data, 0)
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(data)-8))
	pcm, err := DecodeBytes(data, Options{})
	if err != nil {
		t.Fatal(err)
	}
	if pcm.SampleRate != 44100 || len(pcm.Samples) != 3 {
		t.Fatalf("decoded rate=%d samples=%v", pcm.SampleRate, pcm.Samples)
	}
	for i, sample := range pcm.Samples {
		if want := float64(i+1) / (1 << 23); sample != want {
			t.Errorf("sample[%d]=%v, want %v", i, sample, want)
		}
	}
}
