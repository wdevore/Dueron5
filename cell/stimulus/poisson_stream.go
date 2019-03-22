package stimulus

import (
	"math/rand"

	"github.com/wdevore/Deuron5/cell"
	"github.com/wdevore/Deuron5/deuron"
)

func RanGen(seed int64) *rand.Rand {
	return rand.New(rand.NewSource(seed))
}

// PoissonStream generates spikes with poisson distribution (noise).
// The outputs are generally routed into StraitConnections.
type PoissonStream struct {
	basePatternStream

	ran *rand.Rand

	// Random seed
	seed int64

	// The Interspike interval (ISI) counter is populated by a value.
	// When the counter reaches 0 a spike is placed on the output.
	isi int

	// Poisson properties
	firingRate float64
}

// NewPoissonStream creates a stream
func NewPoissonStream(seed int64) IPatternStream {
	s := new(PoissonStream)
	s.baseInitialize()

	s.seed = seed
	s.ran = RanGen(seed)

	s.firingRate = deuron.SimModel.GetFloat("Firing_Rate")

	s.isi = deuron.GenPoisson(s.firingRate)

	return s
}

func (ss *PoissonStream) Initialize(firingRate float64) {
	ss.firingRate = firingRate
	ss.Reset()
}

// Attach connections to this stream
// The given IConnection will have spikes routed into it.
func (ss *PoissonStream) Attach(con cell.IConnection) {
	ss.cons.Add(con)
}

// func (ss *PoissonStream) genPoisson(lambda float64) int {
// 	// return Generate(nps.ran.Float64(), nps.max, nps.spread, nps.min)
// 	return ss.poissonSmall(ss.max)
// }

// func (ss *PoissonStream) poissonSmall(lambda float64) int {
// 	// Algorithm due to Donald Knuth, 1969.
// 	p := 1.0
// 	L := math.Exp(-lambda)
// 	// L := math.Pow(math.E, -lambda)
// 	k := 0
// 	for p > L {
// 		k++
// 		p *= ss.ran.Float64()
// 	}

// 	return k - 1
// }

func (ss *PoissonStream) ISI() int {
	return ss.isi
}

// func Generate(rand, scale, div, min float64) int {
// 	return generate2(rand, scale, div, min)
// }

// func generate1(rand, scale, div, min float64) int {
// 	return deuron.GetPoisson(scale)
// }

// func generate2(rand, scale, div, min float64) int {
// 	return int(scale*math.Pow(math.E, -rand*scale/div) + min)
// }

// func generatePoisson(firingRate float64) int {
// 	return deuron.GenPoisson(firingRate)
// }

// generate tends to spread spikes a bit more.
// Typical values of: 15.0, 3.0 yield ISIs 5-7 with occasional 50-100s,
// or 50.0,15.0,2.0
// func (ss *PoissonStream) generate(scale, div, min float64) int {
// 	return Generate(ss.ran.Float64(), scale, div, min)
// 	// return ss.genPoisson(scale)
// 	// p := deuron.GetPoisson(scale)
// 	// fmt.Printf("%d,", p)
// 	// return p
// }

// ----------------------------------------------
// IPatternStream methods
// ----------------------------------------------
func (ss *PoissonStream) EnableAutoReset() {
	// Not applicable
}

// Reset generates a new ISI
func (ss *PoissonStream) Reset() {
	deuron.SeedPoisson(ss.seed)

	ss.ran.Seed(ss.seed)
	ss.firingRate = deuron.SimModel.GetFloat("Firing_Rate")

	ss.isi = deuron.GenPoisson(ss.firingRate)

	// ss.isi = ss.generate(ss.max, ss.spread, ss.min)
}

func (ss *PoissonStream) Step() bool {

	// Check ISI counter
	if ss.isi == 0 {
		// Time to generate a spike
		ss.value = 1
		// ss.isi = ss.generate(ss.max, ss.spread, ss.min)
		ss.isi = deuron.GenPoisson(ss.firingRate)
	} else {
		ss.value = 0
		ss.isi--
	}

	// Place stream's current output value onto the
	// associated connection(s) input
	it := ss.cons.Iterator()
	for it.Next() {
		conn := it.Value().(cell.IConnection)
		conn.Input(ss.value)
	}

	return false
}

func (ss *PoissonStream) IsComplete() bool {
	return false // This type of stream never completes
}

// ----------------------------------------------
// IBitStream methods
// ----------------------------------------------

func (ss *PoissonStream) Input(v int) {
	// Not applicable.
}

func (ss *PoissonStream) Output() int {
	return ss.value
}
