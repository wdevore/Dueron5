package cell

import (
	"math"
	"strconv"

	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/simulation/samples"
)

const (
	RefractoryEpsilon = 0.1
)

// This neuron is for prototyping only.

type ProtoNeuron struct {
	baseCell

	// Soma threshold. When exceeded an AP is generated.
	threshold float64

	// --------------------------------------------------------
	// Action potential
	// --------------------------------------------------------
	// AP can travel back down the dendrite. The value decays
	// with distance.
	apFast      float64 // Fast trace
	apSlow      float64 // Slow trace
	apSlowPrior float64 // Slow trace (t-1)

	// The time-mark of the current AP.
	APt float64
	// The previous time-mark
	preAPt float64
	APMax  float64

	// --------------------------------------------------------
	// STDP
	// --------------------------------------------------------
	// -----------------------------------
	// AP decay
	// -----------------------------------
	ntao  float64 // fast trace
	ntaoS float64 // slow trace

	// Fast Surge
	nFastSurge        float64
	nDynFastSurge     float64
	nInitialFastSurge float64

	// Slow Surge
	nSlowSurge        float64
	nDynSlowSurge     float64
	nInitialSlowSurge float64

	// The time-mark at which a spike arrived at a synapse
	preT float64

	// -----------------------------------
	// Suppression
	// -----------------------------------
	ntaoJ         float64
	efficacyTrace float64

	// -----------------------------------
	// Fall off
	// -----------------------------------
}

func NewProtoNeuron() ICell {
	n := new(ProtoNeuron)
	n.baseCell.initialize()

	n.Reset()

	return n
}

func (n *ProtoNeuron) Output() float64 {
	return n.output
}

func (n *ProtoNeuron) AddInConnection(con IConnection) {
	n.inputs = append(n.inputs, con)
}

func (n *ProtoNeuron) AddOutConnection(con IConnection) {
	n.outputs = append(n.outputs, con)
}

func (n *ProtoNeuron) Reset() {
	n.apFast = 0.0
	n.apSlow = 0.0
	n.preT = -1000000000000000.0
	n.refractoryState = false
	n.refractoryCnt = 0
	n.nSlowSurge = 0.0
	n.nFastSurge = 0.0
	n.output = 0.0
	n.prevOutput = 0.0
	n.efficacyTrace = 0.0
	if n.dendrite != nil {
		n.dendrite.Reset()
	}
}

func (n *ProtoNeuron) Integrate(t float64) float64 {
	dt := t - n.preT

	samples.Sim.NeuronDtSamples.Put(t, dt, n.id, 0)

	n.efficacyTrace = n.efficacy(dt, n.ntaoJ)

	// Pass the current AP trace and the current and previous spike state
	// contained in the cell.
	psp := n.dendrite.Integrate(t, n)
	samples.Sim.NeuronPspSamples.Put(t, psp, n.id, 0)

	n.prevOutput = n.output

	// Default state
	n.output = 0.0

	if n.refractoryState {
		// this algorithm should be the same as for the synapse or at least very
		// close.

		if n.refractoryCnt >= n.refractoryPeriod {
			n.refractoryState = false
			n.refractoryCnt = 0
			// fmt.Printf("Refractory ended at (%d)\n", int(t))
		} else {
			n.refractoryCnt++
		}
	} else {
		if psp > n.threshold {
			// An action potential just occurred.

			// TODO Handle depolarization

			n.refractoryState = true

			// TODO
			// Generate a back propagating spike that fades spatial/temporally similar to CaDP model.
			// This spike affects forward in time.
			// The value is driven by the time delta of (preAPt - APt)

			n.output = 1.0

			// Surge from action potential
			n.nFastSurge = n.APMax + n.apFast*n.nInitialFastSurge*math.Exp(-n.apFast/n.ntao)
			n.nSlowSurge = n.APMax + n.apSlow*n.nInitialSlowSurge*math.Exp(-n.apSlow/n.ntaoS)

			// Reset time deltas
			n.preT = t
			dt = 0
		}
	}

	// Prior is for triplet
	n.apSlowPrior = n.apSlow

	n.apFast = n.nFastSurge * math.Exp(-dt/n.ntao)
	n.apSlow = n.nSlowSurge * math.Exp(-dt/n.ntaoS)

	samples.Sim.NeuronAPSamples.Put(t, n.apFast, n.id, 0)
	samples.Sim.NeuronAPSlowSamples.Put(t, n.apSlow, n.id, 0)

	return n.output
}

// This is a time based property NOT distance.
// Each spike of the neuron i sets the post spike efficacy j to 0
// whereafter it recovers exponentially to 1 with a time constant toaI.
// In other words, the efficacy of a spike is suppressed by
// the proximity of a previous post spike.
func (n *ProtoNeuron) efficacy(dt, taoI float64) float64 {
	return 1 - math.Exp(-dt/taoI)
}

func (n *ProtoNeuron) Process() {
	n.dendrite.Process()
}

func (n *ProtoNeuron) PostProcess() {
	n.dendrite.PostProcess()
}

// -----------------------------------------------------------------
// Properties
// -----------------------------------------------------------------

func (n *ProtoNeuron) APFast() float64 {
	return n.apFast
}

func (n *ProtoNeuron) APSlow() float64 {
	return n.apSlow
}

func (n *ProtoNeuron) APSlowPrior() float64 {
	return n.apSlowPrior
}

func (n *ProtoNeuron) Efficacy() float64 {
	return n.efficacyTrace
}

func (n *ProtoNeuron) SetThreshold(theta float64) {
	n.threshold = theta
}

func (n *ProtoNeuron) SetField(field, value string) {
	switch field {
	case "nFastSurge":
		n.nInitialFastSurge, _ = strconv.ParseFloat(value, 64)
		break
	case "nSlowSurge":
		n.nInitialSlowSurge, _ = strconv.ParseFloat(value, 64)
		break
	case "ntao":
		n.ntao, _ = strconv.ParseFloat(value, 64)
		break
	case "ntaoS":
		n.ntaoS, _ = strconv.ParseFloat(value, 64)
		break
	case "ntaoJ":
		n.ntaoJ, _ = strconv.ParseFloat(value, 64)
		break
	case "APMax":
		n.APMax, _ = strconv.ParseFloat(value, 64)
		break
	}
}

func (n *ProtoNeuron) Load(json interface{}) {
	jmap := json.(map[string]interface{})

	n.nInitialFastSurge = jmap["nFastSurge"].(float64)
	n.nInitialSlowSurge = jmap["nSlowSurge"].(float64)
	n.ntao = jmap["ntao"].(float64)
	n.ntaoS = jmap["ntaoS"].(float64)
	n.ntaoJ = jmap["ntaoJ"].(float64)

	n.threshold = jmap["Threshold"].(float64)
	n.refractoryPeriod = jmap["RefractoryPeriod"].(float64)
	n.APMax = jmap["APMax"].(float64)

	m := deuron.SimModel

	m.SetFloat("threshold", n.threshold)
	m.SetFloat("RefractoryPeriod", n.refractoryPeriod)
	m.SetFloat("APMax", n.APMax)

	m.SetFloat("nFastSurge", n.nInitialFastSurge)
	m.SetFloat("nSlowSurge", n.nInitialSlowSurge)
	m.SetFloat("ntao", n.ntao)
	m.SetFloat("ntaoS", n.ntaoS)
	m.SetFloat("ntaoJ", n.ntaoJ)

	n.dendrite.Load(json)
}

func (n *ProtoNeuron) ToJSON() interface{} {
	m := map[string]interface{}{
		"id":               n.id,
		"Threshold":        n.threshold,
		"ntao":             n.ntao,
		"ntaoS":            n.ntaoS,
		"ntaoJ":            n.ntaoJ,
		"nFastSurge":       n.nInitialFastSurge,
		"nSlowSurge":       n.nInitialSlowSurge,
		"RefractoryPeriod": n.refractoryPeriod,
		"APMax":            n.APMax,
		"wMin":             deuron.SimModel.GetFloat("weightMin"),
		"wMax":             deuron.SimModel.GetFloat("weightMax"),
		"Dendrites":        n.dendrite.ToJSON(),
	}

	return m
}
