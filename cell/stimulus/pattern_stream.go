package stimulus

import (
	"fmt"
	"strings"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
)

// NPatternStream N spike streams.
// Each stream is routed to 1 or more IConnections.
type NPatternStream struct {
	output byte

	complete  bool
	autoReset bool

	// A collection of streams
	patterns *sll.List
	patItr   sll.Iterator
	// period = 1/frequency * 1000 (ms)
	period int // in milliseconds

	// delay = period - pattern_length
	delay    int
	delayCnt int
}

func NewNPatternStream() *NPatternStream {
	s := new(NPatternStream)
	s.autoReset = true
	s.patterns = sll.New()
	return s
}

// Period sets the "delay" between pattern applications. The pattern
// period itself could be longer that the delay.
func (nps *NPatternStream) Period(period, patternLen int) {
	nps.period = period
	nps.delay = nps.period - patternLen
	nps.delayCnt = nps.delay

}

func (nps *NPatternStream) Add(strm IPatternStream) {
	nps.patterns.Add(strm)
}

func (nps *NPatternStream) IsComplete() bool {
	return false
}

func (nps *NPatternStream) EnableAutoReset() {
	nps.autoReset = true
}

func (nps *NPatternStream) Reset() {
	nps.complete = false
	nps.delayCnt = nps.delay
	it := nps.patterns.Iterator()
	for it.Next() {
		stim := it.Value().(IPatternStream)
		stim.Reset()
	}
}

func (nps *NPatternStream) Step() {
	// Step all the streams when the delay had ended.
	// Once the pattern has completed we switch back to delay.
	if nps.delayCnt < 0 {
		var complete bool
		it := nps.patterns.Iterator()
		for it.Next() {
			stim := it.Value().(IPatternStream)
			complete = complete || stim.Step()
		}

		if complete {
			// Reset for another delay period
			nps.Reset()
		}
	} else {
		nps.delayCnt--
	}
}

func (nps *NPatternStream) Begin() bool {
	nps.patItr = nps.patterns.Iterator()
	return nps.patItr.First()
}

func (nps *NPatternStream) Next() bool {
	return nps.patItr.Next()
}

// Stream returns the next available pattern stream.Begin
func (nps *NPatternStream) Stream() IPatternStream {
	stream := nps.patItr.Value().(IPatternStream)
	return stream
}

func (nps NPatternStream) String() string {
	var s strings.Builder

	it := nps.patterns.Iterator()
	for it.Next() {
		stim := it.Value().(*SpikeStream)
		s.WriteString(fmt.Sprintf("%s\n", stim.String()))
	}

	return s.String()
}
