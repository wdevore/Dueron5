package samples

import (
	"fmt"
	"math"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
)

// ---------------------------------------------------------
// Allocated in run_reset.go Create()
// ---------------------------------------------------------

type baseSamples struct {
	lanes *sll.List

	laneCnt int
	size    int // typically the length of simulation

	// Window range parameters
	RangeStart int
	RangeEnd   int
}

// Samples is a 2D list of samples for synapses
type Samples struct {
	baseSamples
}

type Sample struct {
	Time float64

	// Typically int or float64
	Value interface{}

	// This 'key' can represent anything, for example what color to
	// render the spike on a graph
	Key int

	// What the data belongs to.
	Id int
}

type baseLane struct {
	Id     int
	Values []*Sample
	Min    float64
	Max    float64
}

// Lanes are trains of data for a given synapse or neuron.
type SamplesLane struct {
	baseLane
}

// =======================================================================
// Allocations
// =======================================================================

func NewSamples(synCnt, size int) *Samples {
	s := new(Samples)

	s.lanes = sll.New()
	s.laneCnt = synCnt
	s.size = size

	// Pre expand collection
	for i := 0; i < synCnt; i++ {
		l := new(SamplesLane)
		l.Id = i
		l.Values = make([]*Sample, size)
		for s := 0; s < size; s++ {
			l.Values[s] = new(Sample)
		}
		s.lanes.Add(l)
	}

	s.ResetRange()

	// s.mutex = &sync.Mutex{}
	return s
}

func (s *Samples) GetLanes() *sll.List {
	return s.lanes
}

func (s *Samples) Size() int {
	return s.size
}

func (s *Samples) Print() {
	it := s.lanes.Iterator()
	for it.Next() {
		lane := it.Value().(*SamplesLane)
		fmt.Printf("(%d) %v\n", lane.Id, lane.Values)
	}
}

func (s *Samples) Put(time float64, value interface{}, sid, key int) {
	// sid is usually synId
	_, lif := s.lanes.Find(func(id int, v interface{}) bool {
		return v.(*SamplesLane).Id == sid
	})

	l := lif.(*SamplesLane)
	sp := l.Values[int(time)]

	sp.Time = time
	sp.Value = value
	sp.Id = sid
	sp.Key = key
}

func (s *Samples) Post() {
	it := s.lanes.Iterator()

	// Find min/max values for each lane
	for it.Next() {
		lane := it.Value().(*SamplesLane)

		// Find min/max values for lane
		min := 1000000000000.0
		max := -1000000000000.0

		for i, v := range lane.Values {
			if v.Value == nil {
				panic(fmt.Sprintf("Sample had nil at (%d)\n", i))
			} else {
				min = math.Min(min, v.Value.(float64))
				max = math.Max(max, v.Value.(float64))
			}
		}

		lane.Min = min
		lane.Max = max

		// fmt.Printf("(%d) min: %f, max: %f\n", lane.Id, min, max)
	}
}

// =======================================================================
// Window range functions
// =======================================================================

func (s *Samples) ResetRange() {
	s.RangeStart = 0
	s.RangeEnd = s.size
}

func (s *Samples) SetRangeStart(start int) {
	s.RangeStart = start
}

func (s *Samples) SetRangeEnd(end int) {
	s.RangeEnd = end
}

func (s *Samples) GetRange() (start, end int) {
	return s.RangeStart, s.RangeEnd
}

func (s *Samples) GetRangeWidth() int {
	return s.RangeEnd - s.RangeStart
}

// func (s *Samples) Use() *sll.List {
// 	s.mutex.Lock()
// 	return s.lanes
// }

// func (s *Samples) Release() {
// 	s.mutex.Unlock()
// }
