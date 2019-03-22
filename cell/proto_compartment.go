package cell

type ProtoCompartment struct {
	baseCompartment
}

func NewProtoCompartment(den IDendrite) ICompartment {
	n := new(ProtoCompartment)

	// Bidirectional associations
	n.den = den
	n.den.AddCompartment(n)
	n.baseCompartment.initialize()

	return n
}

func (c *ProtoCompartment) Process() {
	it := c.synapses.Iterator()
	for it.Next() {
		synapse := it.Value().(ISynapse)
		synapse.Process()
	}
}

func (c *ProtoCompartment) PostProcess() {
	it := c.synapses.Iterator()
	for it.Next() {
		synapse := it.Value().(ISynapse)
		synapse.PostProcess()
	}
}

func (c *ProtoCompartment) Dendrite() IDendrite {
	return c.den
}

func (c *ProtoCompartment) Integrate(t float64, cell ICell) float64 {
	// Compute the Post Synaptic Potential (PSP)
	// Note: this PSP may be used later with a diffusion effect.

	psp := 0.0

	it := c.synapses.Iterator()
	for it.Next() {
		synapse := it.Value().(ISynapse)
		psp += synapse.Integrate(t, cell)
	}

	// if psp < 0 {
	// 	psp = 0.0
	// }

	return psp
}

func (c *ProtoCompartment) Load(json interface{}) {
	it := c.synapses.Iterator()

	jmap := json.(map[string]interface{})
	dendrites := jmap["Dendrites"].(map[string]interface{})
	compartments := dendrites["Compartments"].([]interface{})

	compartment := compartments[0].(map[string]interface{})
	synsArr := compartment["Synapses"].([]interface{})

	i := 0
	for it.Next() {
		synapse := it.Value().(ISynapse)
		synapse.Load(synsArr[i])
		i++
	}
}

func (c *ProtoCompartment) ToJSON() interface{} {
	a := make([]interface{}, c.synapses.Size())

	it := c.synapses.Iterator()
	ind := 0
	for it.Next() {
		synapse := it.Value().(ISynapse)
		a[ind] = synapse.ToJSON()
		ind++
	}

	m := map[string]interface{}{
		"id":       c.id,
		"Synapses": a,
	}

	return m
}
