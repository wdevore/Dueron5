package cell

// This type of neuron handles stimulus type inputs, however, it
// can also handle interior/hidden type neurons.
type StimulusNeuron struct {
	baseCell
}

// func NewStimulusNeuron() ICell {
// 	n := new(StimulusNeuron)
// 	n.baseCell.initialize()
// 	return n
// }

// func (n *StimulusNeuron) Output() int {
// 	return n.output
// }

// func (n *StimulusNeuron) AddInConnection(con IConnection) {
// 	n.inputs = append(n.inputs, con)
// }

// func (n *StimulusNeuron) AddOutConnection(con IConnection) {
// 	n.outputs = append(n.outputs, con)
// }

// func (n *StimulusNeuron) Integrate(dt float64) float64 {
// 	return 0.0
// }

// func (n *StimulusNeuron) Process() {

// }

// func (n *StimulusNeuron) Reset() {

// }
