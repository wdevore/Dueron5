package runreset

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"

	"github.com/wdevore/Deuron5/deuron/app/comm"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/wdevore/Deuron5/cell"
	"github.com/wdevore/Deuron5/cell/stimulus"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/simulation/samples"
)

type Simulation struct {
	channel chan string

	neuron cell.ICell

	poiStreams  *sll.List
	stimStreams *sll.List
	syns        *sll.List
	cons        *sll.List

	pattern1 *stimulus.PoissonPatternStream

	cnt int

	settingsMap map[string]interface{}
}

// NewSimulation creates a simulation
func NewSimulation(channel chan string) *Simulation {
	s := new(Simulation)
	s.channel = channel
	return s
}

var ran = rand.New(rand.NewSource(1963))

func (s *Simulation) initialize() int {
	// The single neuron being simulated.
	s.neuron = cell.NewProtoNeuron()

	threshold := deuron.SimModel.GetFloat("threshold")
	s.neuron.SetThreshold(threshold)

	// A neuron has a dendrite
	den := cell.NewProtoDendrite(s.neuron)

	// A dendrite has 1 or more compartments--one in this simulation.
	comp := cell.NewProtoCompartment(den)

	// Create 80% Excite and 20% Inhibit
	// The streams are setup with N channels.
	synCount := int(deuron.SimModel.GetFloat("Synapse_Count"))

	excite := int(float64(synCount) * 0.8)
	inhibit := int(float64(synCount) * 0.2)

	// Collections used for convenience of iteration.
	s.poiStreams = sll.New()
	s.stimStreams = sll.New()
	s.syns = sll.New()
	s.cons = sll.New()

	// Id generators
	synID := 0
	poiID := 0

	s.loadSettings()

	// Now fill streams with patterns
	s.createPatterns()

	s.pattern1.Begin()

	// Begin construction of the neuron.
	// For each synapse we attach a connection.
	// For this simulation each connection is also connected to
	// a poisson-noise stream and pattern stream.
	for i := 0; i < excite; i++ {
		// Create a synapse with a new id and associate it with a compartment and mark it as excititory.
		syn := cell.NewProtoSynapse(comp, cell.Excititory, synID, ran.Int63())

		// Collect it for iteration during simulation.
		s.syns.Add(syn)

		// For this sim we simply connect them without any delays
		con := cell.NewStraightConnection()
		s.cons.Add(con)

		// Create a poisson noise stream that will feed into the connection
		seed := ran.Int63()
		poi := stimulus.NewPoissonStream(seed).(*stimulus.PoissonStream)
		poi.SetId(poiID)

		// Collect streams so we can iterate them later.
		s.poiStreams.Add(poi)

		// -----------------------------------------------------------------
		// Route noise and pattern into connection
		// and connect the pieces together: [stream and noise] -> [connection] -> [synapse]
		// -----------------------------------------------------------------
		poi.Attach(con) // route noise stream into connection

		// Get the stream driving the pattern
		stream := s.pattern1.Stream()
		// Collect it for easy iteration.
		s.stimStreams.Add(stream)

		// route pattern stream into connection as well.
		stream.Attach(con)

		// Finally route connection to synapse
		syn.Connect(con)

		// Generate next ids for debugging.
		synID++
		poiID++
	}

	// Repeat for inhibition.
	for i := 0; i < inhibit; i++ {
		syn := cell.NewProtoSynapse(comp, cell.Inhibitory, synID, ran.Int63())
		s.syns.Add(syn)

		con := cell.NewStraightConnection()
		s.cons.Add(con)

		seed := ran.Int63()
		poi := stimulus.NewPoissonStream(seed).(stimulus.IPatternStream)
		poi.SetId(poiID)

		s.poiStreams.Add(poi)
		// Connect stream to input of connection
		poi.Attach(con)

		stim := s.pattern1.Stream()
		s.stimStreams.Add(stim)
		stim.Attach(con) // route stimulus into connection

		syn.Connect(con) // attach connection into synapse

		synID++
		poiID++
	}

	// Now we can attach dendrite to neuron. Note, neurons can have more than one dendrite.
	s.neuron.AttachDendrite(den)

	model := s.settingsMap["Neuron"]
	s.Load(model)

	fmt.Println("Sim: initialized")

	return synCount
}

// func (s *Simulation) Listen(msg *deuron.MessageEvent) {
// }

func (s *Simulation) reset() {
	it := s.poiStreams.Iterator()
	for it.Next() {
		poi := it.Value().(stimulus.IPatternStream)
		poi.Reset()
	}

	// Reset stimulus
	s.pattern1.Reset()

	// Reset neurons
	s.neuron.Reset()
	s.cnt = 0
}

// A single pass of a simulation.
func (s *Simulation) simulate(t float64) {
	s.pre()

	// Update learning rules (STDP and BTSP) and internal states/properties
	s.neuron.Process()

	// Now integrate
	s.neuron.Integrate(t)

	s.diagnostics(t) // Collect samples for inspection

	s.post()

	// Send message back to App.pollForMessage
	msg := fmt.Sprintf("(%0.1f)", t)
	s.respond(msg)
}

func (s *Simulation) pre() {
	// Step poission noise streams
	it := s.poiStreams.Iterator()
	for it.Next() {
		poi := it.Value().(stimulus.IPatternStream)
		poi.Step()
	}

	// Step all the stimulus streams
	s.pattern1.Step()
}

func (s *Simulation) diagnostics(t float64) {
	// Capture the state at time "t".

	// Collect noise samples from the poisson streams.
	it := s.poiStreams.Iterator()
	for it.Next() {
		pois := it.Value().(stimulus.IPatternStream)
		samples.Sim.PoiSamples.Put(t, pois.Output(), pois.Id(), 3)
	}

	if s.pattern1.Begin() {
		more := true
		for more {
			stim := s.pattern1.Stream()
			if stim == nil {
				more = false
			} else {
				samples.Sim.StimSamples.Put(t, stim.Output(), stim.Id(), 4)
				more = s.pattern1.Next()
			}
		}
	}

	// Capture the cell's current output
	samples.Sim.CellSamples.Put(t, float64(s.neuron.Output()), s.neuron.ID(), 0)
}

func (s *Simulation) Load(json interface{}) {
	// Specific to synapse uniqueness, for example, weights.
	s.neuron.Load(json)
}

// Post process for a single pass
func (s *Simulation) post() {
	// Post is a preperation for next pass.
	// This means we put all Data, on the output side of a connection,
	// back into the pool.

	// The values either source from noise streams, stimulus or other neuron outputs.
	it := s.cons.Iterator()
	for it.Next() {
		con := it.Value().(cell.IConnection)
		con.Post()
	}
}

// Post process for a single simulation run.
func (s *Simulation) PostProcess() {
	s.neuron.PostProcess()

	// Post process any samples.
	// fmt.Println("Post processing...")
	samples.Sim.Post()
}

func (s *Simulation) respond(msg string) {
	// Send message back to the App
	s.channel <- msg
}

// Called from run_reset.go
func (s *Simulation) SendEvent(event *comm.MessageEvent) {
	switch event.Target {
	case "Data":
		switch event.Action {
		case "Changed":
			switch event.Message {
			case "Simulation":
				switch event.Field {
				case "StimulusScaler":
					fValue, _ := strconv.ParseFloat(event.Value, 64)
					s.ExpandStreams(fValue)
					break
				case "Stimulus":
					// Changing stimulus
					expandFactor := int(deuron.SimModel.GetFloat("StimulusScaler"))
					s.loadPatterns(expandFactor)
					s.loadSettings()
					model := s.settingsMap["Neuron"]
					s.Load(model)
					s.respond("GuiRefesh")
					break
				}
				break
			case "Synapse":
				it := s.syns.Iterator()
				for it.Next() {
					synapse := it.Value().(cell.ISynapse)
					synapse.SetField(event.Field, event.Value)
				}
				break
			case "Neuron":
				s.neuron.SetField(event.Field, event.Value)

				threshold := deuron.SimModel.GetFloat("threshold")
				s.neuron.SetThreshold(threshold)
				break
			}
			break
		}
		break
	}
}

func (s *Simulation) createPatterns() {
	// ------------------------------------------------------------
	// Create collection
	s.pattern1 = stimulus.NewPoissonPatternStream(123)
	// s.pattern1.Period(100, 25) // Pattern will be applied at 30Hz or every 33ms

	expandFactor := int(deuron.SimModel.GetFloat("StimulusScaler"))

	s.loadPatterns(expandFactor)
}

func (s *Simulation) loadPatterns(expandFactor int) {
	// Load stimulus patterns
	patFile := "./stimulus/" + deuron.SimModel.GetString("Stimulus") + ".txt"

	patternsFile, err := os.Open(patFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Opened stimulus (%s)\n", patFile)

	defer patternsFile.Close()

	scanner := bufio.NewScanner(patternsFile)
	// spk.SetSpikes([]int{0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 0, 0})

	s.pattern1.Reset()

	if s.stimStreams.Empty() {
		ind := 0
		for scanner.Scan() {
			spk := stimulus.NewSpikeStream().(*stimulus.SpikeStream)

			spk.SetId(ind)
			strm := scanner.Text()
			spk.SetSpikesFromString(strm)
			// fmt.Printf("(%d) %s\n", ind, strm)
			spk.Expand(expandFactor)
			s.pattern1.Add(spk)
			ind++

			s.stimStreams.Add(spk)
		}
	} else {
		it := s.stimStreams.Iterator()

		for scanner.Scan() {
			if it.Next() {
				stream := it.Value().(*stimulus.SpikeStream)
				stream.SetSpikesFromString(scanner.Text())
				stream.Expand(expandFactor)
			}
		}
	}
}

func (s *Simulation) loadSettings() {
	// Load stimulus patterns
	fileName := "./stimulus/" + deuron.SimModel.GetString("Stimulus") + ".json"

	settingsFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Opened settings (%s)\n", fileName)

	defer settingsFile.Close()

	byteValue, _ := ioutil.ReadAll(settingsFile)

	s.settingsMap = make(map[string]interface{})

	err = json.Unmarshal(byteValue, &s.settingsMap)

	if err != nil {
		fmt.Println(err)
		return
	}

	m := deuron.SimModel

	m.SetFloat("StimulusScaler", s.settingsMap["StimulusScaler"].(float64))
	m.SetFloat("Hertz", s.settingsMap["Hertz"].(float64))

	m.SetFloat("Firing_Rate", s.settingsMap["Firing_Rate"].(float64))
	m.SetFloat("Poisson_Pattern_max", s.settingsMap["Poisson_Pattern_max"].(float64))
	m.SetFloat("Poisson_Pattern_spread", s.settingsMap["Poisson_Pattern_spread"].(float64))
	m.SetFloat("Poisson_Pattern_min", s.settingsMap["Poisson_Pattern_min"].(float64))

}

func (s *Simulation) ExpandStreams(scaler float64) {
	s.pattern1.ExpandStreams(scaler)
}

func (s *Simulation) ToJSON() interface{} {
	mo := deuron.SimModel

	m := map[string]interface{}{
		"Firing_Rate":            mo.GetFloat("Firing_Rate"),
		"Poisson_Pattern_max":    mo.GetFloat("Poisson_Pattern_max"),
		"Poisson_Pattern_spread": mo.GetFloat("Poisson_Pattern_spread"),
		"Poisson_Pattern_min":    mo.GetFloat("Poisson_Pattern_min"),

		"threshold":        mo.GetFloat("threshold"),
		"RefractoryPeriod": mo.GetFloat("RefractoryPeriod"),

		"StimulusScaler": mo.GetFloat("StimulusScaler"),
		"Hertz":          mo.GetFloat("Hertz"),

		"Neuron": s.neuron.ToJSON(),
	}

	return m
}
