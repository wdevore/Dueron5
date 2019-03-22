package stimulus

import (
	"fmt"
	"strings"

	"github.com/wdevore/Deuron5/cell"
)

// SpikeStream provides a spiking stimulus stream
// For a pattern to form spatially you will need 2 or more streams
//
// This type of stream outputs a pattern at a certain frequency.
// For example, at 60Hz the pattern is presented every 17ms which means
// the pattern's width needs to be at least < than the period.
//
// ---------------------|--pattern--|----...----|--pattern--|---...-----|--pattern--|
// ^                                ^
// |  period = 17ms                 |
//
type SpikeStream struct {
	basePatternStream

	complete  bool
	autoReset bool

	// Spike pattern. The pattern is fixed in size, for now.
	pattern []int

	// The pattern output and possibly expanded
	expanded []int

	idx int
}

// [size] is in milliseconds
func NewSpikeStream() IPatternStream {
	s := new(SpikeStream)
	s.baseInitialize()

	s.autoReset = false

	s.Reset()
	s.EnableAutoReset()

	return s
}

func (ss *SpikeStream) Input(v int) {
}

func (ss *SpikeStream) Output() int {
	return ss.value
}

func (ss *SpikeStream) IsComplete() bool {
	return ss.complete
}

func (ss *SpikeStream) EnableAutoReset() {
	ss.autoReset = true
}

func (ss *SpikeStream) Reset() {
	ss.complete = false
	// ss.idx = len(ss.expanded) - 1  // end to start
	ss.idx = 0 // start to end
	ss.value = 0
}

func (ss *SpikeStream) Step() bool {
	// Only step the pattern after the delay.

	// Step the pattern from end to start so the
	// pattern in code looks the same on the display.
	if ss.autoReset && ss.idx >= len(ss.expanded) {
		// if ss.autoReset && ss.idx < 0 {
		ss.Reset()
		return true // complete
	}

	ss.value = ss.expanded[ss.idx]

	// Place stream's current output value onto the
	// associated connection(s) input
	it := ss.cons.Iterator()
	for it.Next() {
		conn := it.Value().(cell.IConnection)
		conn.Input(ss.value)
	}

	// ss.idx--
	ss.idx++

	return false // not complete yet
}

// Set sets a particular spike
func (ss *SpikeStream) Set(t int) {
	if t > len(ss.expanded) {
		fmt.Println("SpikeStream: bad t position")
		return
	}

	ss.expanded[t] = 1
}

// SetRange sets a range of spikes from an array to 1
func (ss *SpikeStream) SetRange(ts []int) {
	for _, ti := range ts {

		if ti > len(ss.expanded) {
			fmt.Println("SpikeStream: bad t position")
			return
		}

		ss.expanded[ti] = 1
	}
}

// SetSpikes sets and "loads" a fresh set of spikes.
func (ss *SpikeStream) SetSpikes(sp []int) {
	ss.pattern = make([]int, len(sp))
	ss.expanded = make([]int, len(sp))

	for t, spik := range sp {

		if t > len(ss.pattern) {
			fmt.Println("SpikeStream: bad set position")
			return
		}

		ss.pattern[t] = spik
		ss.expanded[t] = spik
	}
	ss.Reset()
}

func (ss *SpikeStream) SetSpikesFromString(sp string) {
	// s := strings.Replace(sp, " ", "", -1)

	ss.pattern = make([]int, len(sp))
	ss.expanded = make([]int, len(sp))

	t := 0
	for _, spik := range sp {

		if t > len(ss.pattern) {
			fmt.Println("SpikeStream: bad set position")
			return
		}

		if spik == '.' {
			ss.pattern[t] = 0
			ss.expanded[t] = 0
			t++
		} else if spik == '|' {
			ss.pattern[t] = 1
			ss.expanded[t] = 1
			t++
		}
	}

	ss.Reset()
}

func (ss *SpikeStream) Expand(scale int) {
	ss.expanded = []int{}

	for _, spik := range ss.pattern {
		ss.expanded = append(ss.expanded, spik)
		for i := 0; i < scale; i++ {
			ss.expanded = append(ss.expanded, 0)
		}
	}
}

func (ss *SpikeStream) Clear(t int) {
	if t > len(ss.pattern) {
		fmt.Println("SpikeStream: bad clear position")
		return
	}
	ss.pattern[t] = 0
	ss.expanded[t] = 0 // This isn't quite correct "t" should be scaled
}

func (ss SpikeStream) ToString(reverse bool) string {
	var s strings.Builder

	if reverse {
		for j := len(ss.expanded) - 1; j >= 0; j-- {
			if ss.expanded[j] == 0 {
				s.WriteString(".")
			} else {
				s.WriteString("|")
			}
		}
	} else {
		for j := 0; j < len(ss.expanded); j++ {
			if ss.expanded[j] == 0 {
				s.WriteString(".")
			} else {
				s.WriteString("|")
			}
		}
	}

	return s.String()
}

func (ss SpikeStream) String() string {
	var s strings.Builder

	for j := len(ss.expanded) - 1; j >= 0; j-- {
		if ss.expanded[j] == 0 {
			s.WriteString("_ ")
		} else {
			s.WriteString("|")
		}
	}

	return s.String()
}
