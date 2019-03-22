package stimulus

import "github.com/wdevore/Deuron5/cell"

// IBitStream can either be input stimulus or neuron output spikes
// This is a stream of spikes
type IBitStream interface {
	// Output is the current value of the stream
	Output() int

	// Input is an alternate source of data.
	// A stream can either generate output or receive output directly
	// from an Input
	Input(v int)
}

type IPatternStream interface {
	IBitStream

	SetId(id int)

	Id() int

	// This stream will send spikes to the connection.
	Attach(cell.IConnection)

	// IsComplete indicates if the stream has reached the end
	// Note: neurons do NOT have an end so this value is always `false`
	IsComplete() bool

	// EnableAutoReset forces the stream to restart/reset
	// once the stream reaches it end--if it has an "ending".
	EnableAutoReset()

	// Reset stream back to fit first bit. This only applies to
	// stimulus streams.
	Reset()

	// Step moves the stream to it next value
	// returns true if patten complete during this step.
	Step() bool
}
