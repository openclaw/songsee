package audio

import (
	"fmt"
	"math"
)

// Slice returns a time-based slice of audio in seconds.
func Slice(a Audio, startSec, durationSec float64) (Audio, error) {
	if math.IsNaN(startSec) || math.IsInf(startSec, 0) || math.IsNaN(durationSec) || math.IsInf(durationSec, 0) || startSec < 0 || durationSec < 0 {
		return Audio{}, fmt.Errorf("slice: start and duration must be finite and >= 0")
	}
	if a.SampleRate <= 0 {
		return Audio{}, fmt.Errorf("slice: invalid sample rate")
	}
	if len(a.Samples) == 0 {
		return Audio{}, fmt.Errorf("slice: empty samples")
	}

	startSamples := startSec * float64(a.SampleRate)
	if startSamples >= float64(len(a.Samples)) {
		return Audio{}, fmt.Errorf("slice: start beyond end")
	}
	start := int(startSamples)
	end := len(a.Samples)
	if durationSec > 0 {
		// Bound the sample count before conversion or addition can overflow.
		durationSamples := durationSec * float64(a.SampleRate)
		if durationSamples < float64(len(a.Samples)-start) {
			end = start + int(durationSamples)
		}
		if end <= start {
			return Audio{}, fmt.Errorf("slice: duration too short")
		}
	}

	return Audio{SampleRate: a.SampleRate, Samples: a.Samples[start:end]}, nil
}
