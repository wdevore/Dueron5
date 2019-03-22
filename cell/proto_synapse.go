package cell

import (
	"math"
	"strconv"

	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/simulation/samples"
)

// Basic pair-based update rule for STDP
// Potentiate: pre-post
// the update of w at the moment of postsynaptic spike
// is proportional to the momentary value of psp trace.
// When a post (aka AP) occurs we read the synapse trace as the value
// to add to "w"
//
// Depression: post-pre

// This synapse is for prototyping only.
const (
	InitialPreT = 0.0 //-10000000000.0
)

type ProtoSynapse struct {
	baseSynapse

	// Initial values
	jmap map[string]interface{}

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Surge
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Surge base value
	amb float64

	// Surge peak
	ama float64

	// Surge window
	tsw float64

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// new surge ion concentration
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// concentration base. We should always have a minimum concentration
	// as a result of a spike
	// Surge is calculated at the arrival of a spike
	// surge = amb - ama*e^(-psp/tsw) == rising curve
	surge float64

	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	// The time-mark at which a spike arrived at a synapse
	preT float64

	// The current ion concentration
	psp float64

	// =============================================================
	// Learning rules:
	// =============================================================
	//
	// Depression pair-STDP, Potentiation is triplet.
	// "tao"s control the rate of decay. Larger values means a slower decay.
	// Smaller values equals a sharper decay.
	// -----------------------------------

	// denominator, positive window time decay
	taoP float64

	// denominator, negative window time decay
	taoN float64

	// Ratio of mRate/taoX
	tao float64

	// -----------------------------------
	// Weight dependence
	// -----------------------------------
	// F-(w) = ƛ⍺w^µ, F+(w) = ƛ(1-w)^µ
	mu     float64 // µ
	lambda float64 // ƛ
	alpha  float64 // ⍺

	// -----------------------------------
	// Suppression
	// -----------------------------------
	taoI         float64
	prevEffTrace float64

	learningRateSlow float64
	learningRateFast float64

	// -----------------------------------
	// Fall off
	// -----------------------------------
	distanceEfficacy float64
	distance         float64
}

func NewProtoSynapse(comp ICompartment, synType SynapseType, id int, weightSeed int64) ISynapse {
	n := new(ProtoSynapse)
	n.comp = comp
	n.synType = synType
	n.SetId(id)
	comp.AddSynapse(n)

	n.preT = InitialPreT

	n.baseSynapse.initialize()

	// Random weight [Wmin -> Wmax]
	// ran := rand.New(rand.NewSource(weightSeed))
	// n.w = n.wMin + n.wMax*ran.Float64()

	return n
}

func (n *ProtoSynapse) Reset() {
	n.prevEffTrace = 1.0
	n.surge = 0.0
	n.psp = 0
	n.preT = 0
	// Reset weights back to initial values.
	n.wMax = deuron.SimModel.GetFloat("weightMax")
	n.w = n.wMax / 2
}

// It is considered the 1st pass of the simulation per time step.
// 1) Internal values are 'moved' to the outputs.
// 2) Learning rules are applied.
func (n *ProtoSynapse) Process() {

}

// Called at the end of a simulation run (not a pass)
func (n *ProtoSynapse) PostProcess() {
	n.preT = InitialPreT
}

// Integrate is the 2nd pass and handles integration.
// The effects pre/post synaptic spikes are felt here.
func (n *ProtoSynapse) Integrate(t float64, cell ICell) float64 {
	// return n.pair_integrate(t, cell.APFast(), cell.Output())
	return n.triplet_integrate(t, cell)
}

func (n *ProtoSynapse) pair_integrate(t, apFast, somaOutput float64) float64 {
	// Calc psp based on current dynamics: (t - preT). As dt increases
	// psp decreases asymtotically to zero.
	dt := t - n.preT

	samples.Sim.DtSamples.Put(t, dt, n.id, 0)

	// Sample the connection to this synapse. The connection will have already
	// "merged" all traffic through to the connection's output.

	// The output of the connection is the input to this synapse.
	if n.conn.Output() == 1 {
		// Record the time at which the pre-spike arrived for reference in
		// learning rules.
		n.surge = n.amb - n.ama*math.Exp(-n.psp/n.tsw)
		n.psp = n.surge

		// Depression
		// Read post trace and adjust weight accordingly.
		n.w = math.Max(n.w-apFast, n.wMin)

		n.preT = t
		dt = 0.0
	} else {
		// fmt.Printf("%0.3f\n", dt)
		if n.IsExcititory() {
			n.psp = n.surge * math.Exp(-dt/n.taoP)
		} else {
			n.psp = n.surge * math.Exp(-dt/n.taoN)
		}
	}

	samples.Sim.SurgeSamples.Put(t, n.surge, n.id, 0)

	// If an AP occurred we read the current n.psp value and add it
	// to the "w"
	if somaOutput == 1.0 {
		// Potentiation
		// Read pre trace (aka psp) and adjust weight accordingly.
		n.w = math.Min(n.w+n.psp, n.wMax)
	}

	samples.Sim.WeightSamples.Put(t, n.w, n.id, 0)

	// Return the "value" of this synapse for this "t"
	if !n.IsExcititory() {
		samples.Sim.PspSamples.Put(t, -n.psp, n.id, 0)
		return -n.psp * n.w
	}

	samples.Sim.PspSamples.Put(t, n.psp, n.id, 0)

	return n.psp * n.w
}

// Pre trace, Post slow and fast traces.
//
// Depression: fast post trace with at pre spike
// Potentiation: slow post trace at post spike
func (n *ProtoSynapse) triplet_integrate(t float64, cell ICell) float64 {
	// Calc psp based on current dynamics: (t - preT). As dt increases
	// psp decreases asymtotically to zero.
	dt := t - n.preT

	samples.Sim.DtSamples.Put(t, dt, n.id, 0)

	// Sample the connection to this synapse. The connection will have already
	// "merged" all traffic through to the connection's output.

	dwD := 0.0
	dwP := 0.0
	updateWeight := false

	// The output of the axon connection is the input to this synapse.
	if n.conn.Output() == 1 {
		// Record the time at which the pre-spike arrived for reference in
		// learning rules.

		if n.IsExcititory() {
			n.surge = n.psp + n.ama*math.Exp(-n.psp/n.taoP)
		} else {
			n.surge = n.psp + n.ama*math.Exp(-n.psp/n.taoN)
		}

		// #######################################
		// Depression LTD
		// #######################################
		// Read post trace and adjust weight accordingly.
		dwD = n.prevEffTrace * n.weightFactor(false, n.w, n.mu) * cell.APFast()

		n.prevEffTrace = n.efficacy(dt, n.taoI)

		n.preT = t
		dt = 0.0

		updateWeight = true
	}

	if n.IsExcititory() {
		n.psp = n.surge * math.Exp(-dt/n.taoP)
	} else {
		n.psp = n.surge * math.Exp(-dt/n.taoN)
	}

	samples.Sim.SurgeSamples.Put(t, n.surge, n.id, 0)

	// If an AP occurred we read the current n.psp value and add it
	// to the "w"
	if cell.Output() == 1.0 {
		// #######################################
		// Potentiation LTP
		// #######################################
		// Read pre trace (aka psp) and slow AP trace for adjusting weight accordingly.
		//     Post efficacy                                          weight dependence                 triplet sum
		dwP = cell.Efficacy() * n.distanceEfficacy * n.weightFactor(true, n.w, n.mu) * (n.psp + cell.APSlowPrior())
		updateWeight = true
	}

	// Finally update the weight.
	if updateWeight {
		n.w = math.Max(math.Min(n.w+dwP-dwD, n.wMax), n.wMin)
	}

	samples.Sim.WeightSamples.Put(t, n.w, n.id, 0)

	// Return the "value" of this synapse for this "t"
	if !n.IsExcititory() {
		samples.Sim.PspSamples.Put(t, -n.psp, n.id, 0)
		return -n.psp * n.w
	}

	samples.Sim.PspSamples.Put(t, n.psp, n.id, 0)

	return n.psp * n.w
}

// Each spike of pre-synaptic neuron j sets the presynaptic spike
// efficacy j to 0
// whereafter it recovers exponentially to 1 with a time constant
// τj = toaJ
// In other words, the efficacy of a spike is suppressed by
// the proximity of a previous spike.
func (n *ProtoSynapse) efficacy(dt, taoJ float64) float64 {
	return 1 - math.Exp(-dt/taoJ)
}

// mu = 0.0 = additive, mu = 1.0 = multiplicative
func (n *ProtoSynapse) weightFactor(potentiation bool, w, mu float64) float64 {
	if potentiation {
		return n.lambda * math.Pow(1-w/n.wMax, mu)
	}

	return n.lambda * n.alpha * math.Pow(w/n.wMax, mu)
}

// ----------------------------------------------------------------------
// Properties
// ----------------------------------------------------------------------

func (n *ProtoSynapse) SetField(field, value string) {
	switch field {
	case "ama":
		n.ama, _ = strconv.ParseFloat(value, 64)
		break
	case "amb":
		n.amb, _ = strconv.ParseFloat(value, 64)
		break
	case "mu":
		n.mu, _ = strconv.ParseFloat(value, 64)
		break
	case "lambda":
		n.lambda, _ = strconv.ParseFloat(value, 64)
		break
	case "alpha":
		n.alpha, _ = strconv.ParseFloat(value, 64)
		break
	case "learningRateSlow":
		n.learningRateSlow, _ = strconv.ParseFloat(value, 64)
		break
	case "learningRateFast":
		n.learningRateFast, _ = strconv.ParseFloat(value, 64)
		break
	case "taoP":
		n.taoP, _ = strconv.ParseFloat(value, 64)
		break
	case "taoN":
		n.taoN, _ = strconv.ParseFloat(value, 64)
		break
	case "taoI":
		n.taoI, _ = strconv.ParseFloat(value, 64)
		break
	case "distance":
		n.distance, _ = strconv.ParseFloat(value, 64)
		break
	}
}

func (n *ProtoSynapse) SetTaoP(v float64) {
	n.taoP = v
}
func (n *ProtoSynapse) SetTaoN(v float64) {
	n.taoN = v
}
func (n *ProtoSynapse) SetAma(v float64) {
	n.ama = v
}
func (n *ProtoSynapse) SetAmb(v float64) {
	n.amb = v
}

func (n *ProtoSynapse) SetTsw(v float64) {
	n.tsw = v
}

func (n *ProtoSynapse) Load(json interface{}) {
	n.jmap = json.(map[string]interface{})
	n.amb = n.jmap["amb"].(float64)
	n.ama = n.jmap["ama"].(float64)
	n.mu = n.jmap["mu"].(float64)
	n.lambda = n.jmap["lambda"].(float64)
	n.alpha = n.jmap["alpha"].(float64)
	n.learningRateSlow = n.jmap["learningRateSlow"].(float64)
	n.learningRateFast = n.jmap["learningRateFast"].(float64)
	n.taoP = n.jmap["taoP"].(float64)
	n.taoN = n.jmap["taoN"].(float64)
	n.taoI = n.jmap["taoI"].(float64)
	n.distance = n.jmap["distance"].(float64)

	// Calc this synapses's reaction to the AP based on its
	// distance from the soma.
	n.distanceEfficacy = n.comp.Dendrite().APEfficacy(n.distance)

	n.w = n.jmap["w"].(float64)

	n.wMax = deuron.SimModel.GetFloat("weightMax")
	n.wMin = deuron.SimModel.GetFloat("weightMin")

	m := deuron.SimModel

	m.SetFloat("amb", n.amb)
	m.SetFloat("ama", n.ama)
	m.SetFloat("mu", n.mu)
	m.SetFloat("lambda", n.lambda)
	m.SetFloat("alpha", n.alpha)
	m.SetFloat("learningRateSlow", n.learningRateSlow)
	m.SetFloat("learningRateFast", n.learningRateFast)
	m.SetFloat("taoP", n.taoP)
	m.SetFloat("taoN", n.taoN)
	m.SetFloat("taoI", n.taoI)
	m.SetFloat("distance", n.distance)
}

func (n *ProtoSynapse) ToJSON() interface{} {
	m := map[string]interface{}{
		"id":               n.id,
		"w":                n.w,
		"taoP":             n.taoP,
		"taoN":             n.taoN,
		"taoI":             n.taoI,
		"distance":         n.distance,
		"ama":              n.ama,
		"amb":              n.amb,
		"mu":               n.mu,
		"lambda":           n.lambda,
		"alpha":            n.alpha,
		"learningRateSlow": n.learningRateSlow,
		"learningRateFast": n.learningRateFast,
	}

	return m
}
