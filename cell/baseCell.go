package cell

import "github.com/wdevore/Deuron5/deuron"

type baseCell struct {
	id int

	output     float64 // 1 or 0
	prevOutput float64

	// A cell will read its input connections on each pass
	inputs  []IConnection
	outputs []IConnection

	// TODO change to List for more complex sims.
	dendrite IDendrite

	refractoryPeriod float64
	refractoryCnt    float64
	refractoryState  bool
}

func (bc *baseCell) initialize() {
	bc.inputs = []IConnection{}
	bc.outputs = []IConnection{}
	bc.refractoryPeriod = deuron.SimModel.GetFloat("RefractoryPeriod")
}

func (bc *baseCell) AttachDendrite(den IDendrite) {
	bc.dendrite = den
}

func (bc *baseCell) Diagnostics(msg string) {
}

func (bc *baseCell) ID() int {
	return bc.id
}

func (bc *baseCell) SetID(id int) {
	bc.id = id
}

func (bc *baseCell) Load(file string) {

}

func (bc *baseCell) Store(file string) {

}
