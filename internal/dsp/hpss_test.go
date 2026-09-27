package dsp

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
	"testing"
)

func TestMedian9(t *testing.T) {
	// Every permutation covers distinct values; binary inputs cover all duplicate patterns.
	v := []float64{-4, -3, -2, -1, 0, 1, 2, 3, 4}
	var permute func(int)
	permute = func(i int) {
		if i == len(v) {
			if got := median9(v); got != 0 {
				t.Fatalf("median9(%v) = %v", v, got)
			}
			return
		}
		for j := i; j < len(v); j++ {
			v[i], v[j] = v[j], v[i]
			permute(i + 1)
			v[i], v[j] = v[j], v[i]
		}
	}
	permute(0)
	for mask := 0; mask < 1<<9; mask++ {
		for i := range v {
			v[i] = float64(mask >> i & 1)
		}
		if got, want := median9(v), medianOriginal(v); got != want {
			t.Fatalf("median9(%v) = %v, want %v", v, got, want)
		}
	}
}

func TestHPSSBitExact(t *testing.T) {
	rng := rand.New(rand.NewPCG(123, 456))
	// Eight float64 lanes is the maximum supported width (512 bits).
	const maxLanes = 8
	for frames := 0; frames <= 3*maxLanes+7; frames++ {
		for bins := 0; bins <= 3*maxLanes+7; bins++ {
			spec := Spectrogram{Frames: frames, Bins: bins, Values: make([]float64, frames*bins)}
			for i := range spec.Values {
				spec.Values[i] = rng.Float64()*180 - 120
			}
			checkHPSSExact(t, &spec, 9, 9)
		}
	}
	for _, dims := range [][2]int{{1, 1025}, {8, 1025}, {9, 1025}, {33, 1025}, {65, 257}} {
		for _, widths := range [][2]int{{0, 0}, {1, 1}, {2, 4}, {5, 7}, {8, 9}, {9, 8}, {9, 9}, {11, 13}, {40, 38}} {
			t.Run(fmt.Sprintf("%dx%d/%dx%d", dims[0], dims[1], widths[0], widths[1]), func(t *testing.T) {
				spec := Spectrogram{Frames: dims[0], Bins: dims[1], Values: make([]float64, dims[0]*dims[1])}
				for i := range spec.Values {
					spec.Values[i] = rng.Float64()*180 - 120
				}
				checkHPSSExact(t, &spec, widths[0], widths[1])
				values := []float64{math.Inf(-1), -120, -80, -80, math.Copysign(0, -1), 0, 0, 10, math.Inf(1), math.MaxFloat64, -math.MaxFloat64, math.SmallestNonzeroFloat64, -math.SmallestNonzeroFloat64}
				for i := range spec.Values {
					spec.Values[i] = values[rng.IntN(len(values))]
				}
				checkHPSSExact(t, &spec, widths[0], widths[1])
				if len(spec.Values) > 0 {
					spec.Values[len(spec.Values)/2] = math.NaN()
					checkHPSSExact(t, &spec, widths[0], widths[1])
				}
			})
		}
	}
	for _, v := range []float64{0, math.Copysign(0, -1), math.Inf(-1), math.Inf(1), math.NaN()} {
		spec := Spectrogram{Frames: 17, Bins: 31, Values: make([]float64, 17*31)}
		for i := range spec.Values {
			spec.Values[i] = v
		}
		checkHPSSExact(t, &spec, 9, 9)
	}
	t.Log("all HPSS values and Min/Max metadata match origin/main bit-for-bit; maximum deviation 0 ULP")
}

func checkHPSSExact(t *testing.T, spec *Spectrogram, timeWidth, freqWidth int) {
	t.Helper()
	before := append([]float64(nil), spec.Values...)
	wantH, wantP := hpssOriginal(spec, timeWidth, freqWidth)
	gotH, gotP := HPSS(spec, timeWidth, freqWidth)
	for i, v := range spec.Values {
		if math.Float64bits(v) != math.Float64bits(before[i]) {
			t.Fatalf("HPSS mutated input at %d", i)
		}
	}
	for i, pair := range [][2]FeatureMap{{gotH, wantH}, {gotP, wantP}} {
		got, want := pair[0], pair[1]
		if got.Width != want.Width || got.Height != want.Height {
			t.Fatalf("map %d dimensions differ", i)
		}
		assertBits := func(index int, a, b float64) {
			if math.Float64bits(a) != math.Float64bits(b) {
				t.Fatalf("%dx%d widths %dx%d map %d index %d: %x (%v), want %x (%v)", spec.Frames, spec.Bins, timeWidth, freqWidth, i, index, math.Float64bits(a), a, math.Float64bits(b), b)
			}
		}
		assertBits(-2, got.Min, want.Min)
		assertBits(-1, got.Max, want.Max)
		for j := range got.Values {
			assertBits(j, got.Values[j], want.Values[j])
		}
	}
}

func TestHPSSAllocations(t *testing.T) {
	spec := Spectrogram{Frames: 33, Bins: 65, Values: make([]float64, 33*65)}
	for i := range spec.Values {
		spec.Values[i] = float64(i%13) - 10
	}
	if allocations := testing.AllocsPerRun(5, func() { HPSS(&spec, 9, 9) }); allocations > 16 {
		t.Fatalf("HPSS allocated %v times; window processing must reuse scratch buffers", allocations)
	}
}

func TestPowerToDBFloor(t *testing.T) {
	if got := powerToDB(0); got != -120 {
		t.Fatalf("powerToDB(0) = %v, want finite -120", got)
	}
}

func medianOriginal(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	tmp := append([]float64(nil), values...)
	sort.Float64s(tmp)
	mid := len(tmp) / 2
	if len(tmp)%2 == 0 {
		return 0.5 * (tmp[mid-1] + tmp[mid])
	}
	return tmp[mid]
}

// hpssOriginal is the origin/main algorithm, including copying and sorting each window.
func hpssOriginal(spec *Spectrogram, timeWidth, freqWidth int) (harm, perc FeatureMap) {
	if timeWidth <= 0 {
		timeWidth = 9
	}
	if freqWidth <= 0 {
		freqWidth = 9
	}
	frames := spec.Frames
	bins := spec.Bins
	harm = NewFeatureMap(frames, bins)
	perc = NewFeatureMap(frames, bins)

	timeRadius := timeWidth / 2
	freqRadius := freqWidth / 2

	timeBuf := make([]float64, 0, timeWidth)
	freqBuf := make([]float64, 0, freqWidth)
	for f := 0; f < frames; f++ {
		for b := 0; b < bins; b++ {
			timeBuf = timeBuf[:0]
			for tf := f - timeRadius; tf <= f+timeRadius; tf++ {
				if tf < 0 || tf >= frames {
					continue
				}
				timeBuf = append(timeBuf, spec.Values[tf*bins+b])
			}
			freqBuf = freqBuf[:0]
			for tb := b - freqRadius; tb <= b+freqRadius; tb++ {
				if tb < 0 || tb >= bins {
					continue
				}
				freqBuf = append(freqBuf, spec.Values[f*bins+tb])
			}
			hMed := medianOriginal(timeBuf)
			pMed := medianOriginal(freqBuf)
			hPow := dbToPower(hMed)
			pPow := dbToPower(pMed)
			src := dbToPower(spec.Values[f*bins+b])
			den := hPow + pPow + 1e-12
			hVal := src * hPow / den
			pVal := src * pPow / den
			harm.Set(f, b, powerToDB(hVal))
			perc.Set(f, b, powerToDB(pVal))
		}
	}
	return harm, perc
}
