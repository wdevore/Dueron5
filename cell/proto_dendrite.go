package cell

import (
	"math"
	"strconv"

	"github.com/wdevore/Deuron5/deuron"
)

type ProtoDendrite struct {
	baseDendrite

	neuron ICell
}

func NewProtoDendrite(cell ICell) IDendrite {
	n := new(ProtoDendrite)

	// Bidirectional associations
	n.neuron = cell
	n.baseDendrite.initialize()

	return n
}

// 1st pass

// Process handles post processing before Integration is performed.
func (d *ProtoDendrite) Process() {
	it := d.compartments.Iterator()
	for it.Next() {
		comp := it.Value().(ICompartment)
		comp.Process()
	}
}

func (d *ProtoDendrite) PostProcess() {
	it := d.compartments.Iterator()
	for it.Next() {
		comp := it.Value().(ICompartment)
		comp.PostProcess()
	}
}

func (d *ProtoDendrite) APEfficacy(distance float64) float64 {
	if distance < d.length {
		return 1.0
	}

	return math.Exp(-(d.length - distance) / d.taoEff)
}

// 2nd pass
// Integrate is the 2nd pass performing integration.
func (d *ProtoDendrite) Integrate(t float64, cell ICell) float64 {
	psp := 0.0

	it := d.compartments.Iterator()
	for it.Next() {
		comp := it.Value().(ICompartment)
		psp += comp.Integrate(t, cell)
	}

	if psp < 0 {
		psp = 0.0
	}

	return psp
}

func (d *ProtoDendrite) SetField(field, value string) {
	switch field {
	case "length":
		d.length, _ = strconv.ParseFloat(value, 64)
		break
	case "taoEff":
		d.taoEff, _ = strconv.ParseFloat(value, 64)
		break
	}
}

func (d *ProtoDendrite) Load(json interface{}) {
	jmap := json.(map[string]interface{})
	dendrites := jmap["Dendrites"].(map[string]interface{})

	d.length = dendrites["length"].(float64)
	d.taoEff = dendrites["taoEff"].(float64)

	m := deuron.SimModel

	m.SetFloat("length", d.length)
	m.SetFloat("taoEff", d.taoEff)

	it := d.compartments.Iterator()
	for it.Next() {
		comp := it.Value().(ICompartment)
		comp.Load(json)
	}
}

func (d *ProtoDendrite) ToJSON() interface{} {
	a := make([]interface{}, d.compartments.Size())

	it := d.compartments.Iterator()
	ind := 0
	for it.Next() {
		comp := it.Value().(ICompartment)
		a[ind] = comp.ToJSON()
		ind++
	}

	m := map[string]interface{}{
		"id":           d.id,
		"length":       d.length,
		"taoEff":       d.taoEff,
		"Compartments": a,
	}

	return m
}
