package stimulus

import (
	"fmt"
	"math"
	"math/rand"
	"strings"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/wdevore/Deuron5/deuron"
)

// PoissonPatternStream emits a pattern inbetween ISI windows.
// When the pattern has completed emission a new ISI is generated
// as a delay.
type PoissonPatternStream struct {
	output byte

	ran  *rand.Rand
	seed int64

	// Poisson properties
	max    float64
	spread float64
	min    float64

	autoReset bool
	// When a pattern completes generate a new ISI automatically.
	autoGenerateISI bool

	// A collection of streams
	patterns *sll.List
	patItr   sll.Iterator
	// isi = spike intervals
	isi int // in milliseconds

	delayCnt int
}

func NewPoissonPatternStream(seed int64) *PoissonPatternStream {
	s := new(PoissonPatternStream)
	s.autoReset = true
	s.seed = seed
	s.ran = RanGen(seed)

	// Query model for initial values.
	s.max = deuron.SimModel.GetFloat("Poisson_Pattern_max")
	s.spread = deuron.SimModel.GetFloat("Poisson_Pattern_spread")
	s.min = deuron.SimModel.GetFloat("Poisson_Pattern_min")

	// s.isi = Generate(s.ran.Float64(), s.max, s.spread, s.min)
	s.isi = s.genPoisson(s.max)

	s.patterns = sll.New()
	return s
}

func (nps *PoissonPatternStream) genPoisson(lambda float64) int {
	// return Generate(nps.ran.Float64(), nps.max, nps.spread, nps.min)
	return nps.poissonSmall(nps.max)
}

func (nps *PoissonPatternStream) poissonSmall(lambda float64) int {
	// Algorithm due to Donald Knuth, 1969.
	p := 1.0
	L := math.Exp(-lambda)
	// L := math.Pow(math.E, -lambda)
	k := 0
	for p > L {
		k++
		p *= nps.ran.Float64()
	}

	return k - 1
}

// SetISI sets the "ISI" between pattern applications. The pattern
// period itself could be longer than the interval.
func (nps *PoissonPatternStream) SetISI(isi int) {
	nps.isi = isi
}

func (nps *PoissonPatternStream) Add(strm IPatternStream) {
	nps.patterns.Add(strm)
}

func (nps *PoissonPatternStream) IsComplete() bool {
	return false
}

func (nps *PoissonPatternStream) EnableAutoReset() {
	nps.autoReset = true
}

func (nps *PoissonPatternStream) Reset() {
	// fmt.Println("--------------- POI pattern RESETing")
	nps.ran.Seed(nps.seed)
	deuron.SetPoissonSeed(1963)
	nps.patternReset()
}

func (nps *PoissonPatternStream) Clear() {
	nps.patterns.Clear()
}

func (nps *PoissonPatternStream) patternReset() {
	nps.delayCnt = 0

	nps.isi = int(deuron.SimModel.GetFloat("Hertz"))
	if nps.isi == 0.0 {
		nps.max = deuron.SimModel.GetFloat("Poisson_Pattern_max")
		nps.spread = deuron.SimModel.GetFloat("Poisson_Pattern_spread")
		nps.min = deuron.SimModel.GetFloat("Poisson_Pattern_min")

		// nps.isi = Generate(nps.ran.Float64(), nps.max, nps.spread, nps.min)
		nps.isi = nps.poissonSmall(nps.max)
	} else {
		// Convert hertz to isi period
		nps.isi = 1000.0 / nps.isi
	}

	// fmt.Printf("isi: %d\n", nps.isi)
	it := nps.patterns.Iterator()
	for it.Next() {
		stim := it.Value().(IPatternStream)
		stim.Reset()
	}
}

func (nps *PoissonPatternStream) Step() {
	// Step all the streams when the ISI had ended.
	// Once the pattern has completed we switch back to ISI.
	if nps.delayCnt > nps.isi {
		var complete bool
		it := nps.patterns.Iterator()
		for it.Next() {
			stim := it.Value().(IPatternStream)
			complete = complete || stim.Step()
		}

		if complete {
			nps.patternReset()
		}
	} else {
		nps.delayCnt++
	}
}

func (nps *PoissonPatternStream) Begin() bool {
	if nps.patterns.Empty() {
		return false
	}

	nps.patItr = nps.patterns.Iterator()
	nps.patItr.Begin()
	return nps.patItr.Next()
}

func (nps *PoissonPatternStream) Next() bool {
	if nps.patterns.Empty() {
		return false
	}

	return nps.patItr.Next()
}

// Stream returns the current pattern available
func (nps *PoissonPatternStream) Stream() IPatternStream {
	if nps.patterns.Empty() {
		return nil
	}
	it := nps.patItr
	stream := it.Value().(IPatternStream)
	return stream
}

func (nps *PoissonPatternStream) ExpandStreams(scaler float64) {
	if nps.patterns.Empty() {
		return
	}

	it := nps.patterns.Iterator()
	for it.Next() {
		stim := it.Value().(*SpikeStream)
		stim.Expand(int(scaler))
	}
}

func (nps PoissonPatternStream) String() string {
	var s strings.Builder

	it := nps.patterns.Iterator()
	for it.Next() {
		stim := it.Value().(*SpikeStream)
		s.WriteString(fmt.Sprintf("%s\n", stim.String()))
	}

	return s.String()
}
