package cell

import sll "github.com/emirpasic/gods/lists/singlylinkedlist"

// IDendrite collects and manages ICompartments.
type IDendrite interface {
	AddCompartment(ICompartment)

	Integrate(t float64, cell ICell) float64

	Process()

	PostProcess()

	APEfficacy(distance float64) float64

	SetField(field string, value string)

	Reset()

	Load(json interface{})

	ToJSON() interface{}
}

type baseDendrite struct {
	id int

	// Collection of dendrite compartments.
	// Proximal, Apical and Distal
	compartments *sll.List

	length float64
	taoEff float64
}

func (bc *baseDendrite) initialize() {
	bc.compartments = sll.New()

}

func (bc *baseDendrite) Reset() {
	// Reset properties
	it := bc.compartments.Iterator()
	for it.Next() {
		comp := it.Value().(ICompartment)
		comp.Reset()
	}
}

func (bc *baseDendrite) AddCompartment(comp ICompartment) {
	bc.compartments.Add(comp)
}
