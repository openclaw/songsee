package dsp

import "math"

// HPSS separates harmonic and percussive content using median filters.
func HPSS(spec *Spectrogram, timeWidth, freqWidth int) (harm, perc FeatureMap) {
	if timeWidth <= 0 {
		timeWidth = 9
	}
	if freqWidth <= 0 {
		freqWidth = 9
	}
	frames, bins := spec.Frames, spec.Bins
	harm = NewFeatureMap(frames, bins)
	perc = NewFeatureMap(frames, bins)
	timeRadius, freqRadius := timeWidth/2, freqWidth/2
	timeBuf := make([]float64, 0, 2*timeRadius+1)
	freqBuf := make([]float64, 0, 2*freqRadius+1)
	hMedians := make([]float64, bins)
	pMedians := make([]float64, bins)
	// Min/Max propagate NaNs; sorting instead places them before finite values.
	networkOK := true
	for _, v := range spec.Values {
		if math.IsNaN(v) {
			networkOK = false
			break
		}
	}
	const tileFrames = 32
	hTile := make([]float64, min(frames, tileFrames)*bins)
	pTile := make([]float64, len(hTile))
	for first := 0; first < frames; first += tileFrames {
		count := min(tileFrames, frames-first)
		for offset := 0; offset < count; offset++ {
			f := first + offset
			hpssMedians(spec, f, timeRadius, freqRadius, timeBuf, freqBuf, hMedians, pMedians, networkOK)
			for b := 0; b < bins; b++ {
				hPow := dbToPower(hMedians[b])
				pPow := dbToPower(pMedians[b])
				src := dbToPower(spec.Values[f*bins+b])
				den := hPow + pPow + 1e-12
				hVal := src * hPow / den
				pVal := src * pPow / den
				hDB, pDB := powerToDB(hVal), powerToDB(pVal)
				hTile[offset*bins+b], pTile[offset*bins+b] = hDB, pDB
				// Preserve the original frame-first order for the extrema, including ties.
				if hDB < harm.Min {
					harm.Min = hDB
				}
				if hDB > harm.Max {
					harm.Max = hDB
				}
				if pDB < perc.Min {
					perc.Min = pDB
				}
				if pDB > perc.Max {
					perc.Max = pDB
				}
			}
		}
		// FeatureMap is bin-major; flush contiguous frame tiles instead of strided writes.
		for b := 0; b < bins; b++ {
			base := b*frames + first
			for offset := 0; offset < count; offset++ {
				harm.Values[base+offset] = hTile[offset*bins+b]
				perc.Values[base+offset] = pTile[offset*bins+b]
			}
		}
	}
	return harm, perc
}

func hpssMedians(spec *Spectrogram, f, timeRadius, freqRadius int, timeBuf, freqBuf, harm, perc []float64, networkOK bool) {
	for b := 0; b < spec.Bins; b++ {
		harm[b] = hpssTimeMedian(spec, f, b, timeRadius, timeBuf, networkOK)
		perc[b] = hpssFreqMedian(spec, f, b, freqRadius, freqBuf, networkOK)
	}
}

func hpssTimeMedian(spec *Spectrogram, f, b, radius int, buf []float64, networkOK bool) float64 {
	buf = buf[:0]
	for tf := max(0, f-radius); tf <= min(spec.Frames-1, f+radius); tf++ {
		buf = append(buf, spec.Values[tf*spec.Bins+b])
	}
	return hpssMedian(buf, networkOK && radius == 4)
}

func hpssFreqMedian(spec *Spectrogram, f, b, radius int, buf []float64, networkOK bool) float64 {
	buf = append(buf[:0], spec.Values[f*spec.Bins+max(0, b-radius):f*spec.Bins+min(spec.Bins, b+radius+1)]...)
	return hpssMedian(buf, networkOK && radius == 4)
}

func hpssMedian(values []float64, networkOK bool) float64 {
	if len(values) == 9 && networkOK {
		return median9(values)
	}
	return median(values)
}

func median9(v []float64) float64 {
	p0, p1, p2 := v[0], v[1], v[2]
	p3, p4, p5 := v[3], v[4], v[5]
	p6, p7, p8 := v[6], v[7], v[8]
	// The 19-comparator median network, omitting unused outputs.
	p1, p2 = min(p1, p2), max(p1, p2)
	p4, p5 = min(p4, p5), max(p4, p5)
	p7, p8 = min(p7, p8), max(p7, p8)
	p0, p1 = min(p0, p1), max(p0, p1)
	p3, p4 = min(p3, p4), max(p3, p4)
	p6, p7 = min(p6, p7), max(p6, p7)
	p1, p2 = min(p1, p2), max(p1, p2)
	p4, p5 = min(p4, p5), max(p4, p5)
	p7, p8 = min(p7, p8), max(p7, p8)
	p3 = max(p0, p3)
	p5 = min(p5, p8)
	p4, p7 = min(p4, p7), max(p4, p7)
	p6 = max(p3, p6)
	p4 = max(p1, p4)
	p2 = min(p2, p5)
	p4 = min(p4, p7)
	p4, p2 = min(p4, p2), max(p4, p2)
	p4 = max(p6, p4)
	p4 = min(p4, p2)
	return p4
}
