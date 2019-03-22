package samples

import sll "github.com/emirpasic/gods/lists/singlylinkedlist"

// ---------------------------------------------------------
// Allocated in run_reset.go Create()
// ---------------------------------------------------------

var Sim *SamplesCollection

type SamplesCollection struct {
	PoiSamples  *Samples
	StimSamples *Samples

	// Holds surge values from synapses
	SurgeSamples        *Samples
	PspSamples          *Samples
	NeuronPspSamples    *Samples // only one lane
	NeuronAPSamples     *Samples // only one lane
	NeuronAPSlowSamples *Samples // only one lane
	WeightSamples       *Samples

	DtSamples       *Samples
	NeuronDtSamples *Samples

	CellSamples *Samples

	// Collect all the samples that need post processing
	postSamples *sll.List
}

func NewSamplesCollection(synCnt, size int) *SamplesCollection {
	sc := new(SamplesCollection)
	sc.PoiSamples = NewSamples(synCnt, size)
	sc.StimSamples = NewSamples(synCnt, size)
	sc.CellSamples = NewSamples(1, size)

	sc.postSamples = sll.New()

	sc.SurgeSamples = NewSamples(synCnt, size)
	sc.postSamples.Add(sc.SurgeSamples)

	sc.PspSamples = NewSamples(synCnt, size)
	sc.postSamples.Add(sc.PspSamples)

	sc.NeuronPspSamples = NewSamples(1, size)
	sc.postSamples.Add(sc.NeuronPspSamples)

	sc.NeuronAPSamples = NewSamples(1, size)
	sc.postSamples.Add(sc.NeuronAPSamples)
	sc.NeuronAPSlowSamples = NewSamples(1, size)
	sc.postSamples.Add(sc.NeuronAPSlowSamples)

	sc.WeightSamples = NewSamples(synCnt, size)
	sc.postSamples.Add(sc.WeightSamples)

	sc.DtSamples = NewSamples(synCnt, size)
	sc.postSamples.Add(sc.DtSamples)
	sc.NeuronDtSamples = NewSamples(1, size)
	sc.postSamples.Add(sc.NeuronDtSamples)

	return sc
}

func (sc *SamplesCollection) Post() {
	it := sc.postSamples.Iterator()
	for it.Next() {
		s := it.Value().(*Samples)
		s.Post()
	}
}
