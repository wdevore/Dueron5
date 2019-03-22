package cell

import (
	astk "github.com/emirpasic/gods/stacks/arraystack"
)

// Note: This Data concept isn't used at the moment.
// The intent is to track who generated the data and at what time.
// IConnections will "wipeout" the "who" but not the time value.

var PoolData *DataPool = NewDataPool()

type SourceType int

const (
	StimulusSource SourceType = iota
	PoissonSource
	NeuronSource
	// Typically indicates multiple sources spiked at the same time.
	MultiSource
)

// Data represents information flowing through the network.
type Data struct {
	// The time mark at which the data began flowing.
	Time float64

	// Value can be 1 or 0
	Value byte

	// Source indicates where the spike came from.
	Source SourceType
}

func NewData() *Data {
	d := new(Data)
	return d
}

// --------------------------------------------------------
// Data pool based on a stack
// --------------------------------------------------------

type DataPool struct {
	stackP *astk.Stack

	chunkAllocSize int
}

func NewDataPool() *DataPool {
	p := new(DataPool)
	p.stackP = astk.New()
	p.chunkAllocSize = 100
	return p
}

func (dp *DataPool) Get() *Data {
	d, _ := dp.stackP.Pop()

	if d == nil {
		// Put a new chunk of items in the pool and return
		// one of them.
		for i := 0; i < dp.chunkAllocSize; i++ {
			d := NewData()
			dp.stackP.Push(d)
		}
		d, _ = dp.stackP.Pop()
	}

	return d.(*Data)
}

func (dp *DataPool) Put(d *Data) {
	dp.stackP.Push(d)
}

func (dp *DataPool) Drop() {
	dp.stackP.Clear()
}
