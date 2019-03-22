package cell

import sll "github.com/emirpasic/gods/lists/singlylinkedlist"

// ICompartment collects synapses. It represents the functionality
// of a group of synapses located along the Dendrite.
// Compartments are either close to or farther away from the neuron's
// soma: Proximal, Apical, Distal.
// TODO compartments can overlap causing effects to diffuse into neighboring
// compartments, for example, Ca++ dynamics (CaDP).
type ICompartment interface {
	// Compartment properties
	AddSynapse(ISynapse)

	// Behaviors

	// Evaluates the total effective weight for the compartment.
	Integrate(t float64, cell ICell) float64

	Process()

	PostProcess()

	Dendrite() IDendrite

	Reset()

	Load(json interface{})

	ToJSON() interface{}
}

type baseCompartment struct {
	id int

	// Collection of synapses
	synapses *sll.List

	den IDendrite
}

func (bc *baseCompartment) initialize() {
	bc.synapses = sll.New()
}

func (bc *baseCompartment) Reset() {
	it := bc.synapses.Iterator()
	for it.Next() {
		synapse := it.Value().(ISynapse)
		synapse.Reset()
	}
}

func (bc *baseCompartment) AddSynapse(syn ISynapse) {
	bc.synapses.Add(syn)
}
