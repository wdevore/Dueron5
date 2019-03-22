package deuron

import (
	"fmt"
	"strconv"
	"sync"

	hmap "github.com/emirpasic/gods/maps/hashmap"
	"github.com/wdevore/Deuron5/deuron/app/comm"
)

var SimModel = NewModel()

// Model contains all data for the simulation
// The model listens for bus messages and updates properties.
type Model struct {
	props *hmap.Map

	mapMutex *sync.Mutex
}

func NewModel() *Model {
	m := new(Model)
	m.props = hmap.New()

	m.mapMutex = &sync.Mutex{}

	m.initialize()
	return m
}

func (m *Model) initialize() {
	m.props.Put("Status", "Idle")

	// 1 second = 1000ms = 1000000us
	//          = 0.001s  = 0.000001s
	// 10us = 0.00001s = 10/1000000

	// Time step units. Could be 1ms or 100us = 0.1ms or 10us = 0.01ms
	// Convert microseconds to seconds: microseconds/1000000

	m.props.Put("TimeStep", 0.0) // in microseconds

	// Samples length = Duration (seconds) * 1000000 / Timestep (microseconds)
	// Example:
	// 2 second duration at 10us steps = 2*1000000/10 = 200000 samples.
	// Example:
	// 2 second duration at 1ms(1000us) steps = 2*1000000/1000 = 2000 samples.

	// Duration length also equal to sample length.
	// Accessed in App.start()
	m.props.Put("Duration", 0.0)
	m.props.Put("Samples", 0.0)

	m.props.Put("AutoRunPause", 0.0) // 0 = false, 1 = true
	m.props.Put("RangeSync", 0.0)    // 0 = false, 1 = true

	m.props.Put("Inc/Dec", 1.0)

	// ###############################################################
	// BEGIN
	// When one of these properties change then the gui is updated
	// ###############################################################
	// Shrinks or expands ISI for stimulus
	m.props.Put("StimulusScaler", .0)
	// Stimulus file
	m.props.Put("Stimulus", "")

	// If Hertz = 0 then stimulus is distributed as poisson.
	// Hertz is = cycles per second (or 1000ms per second)
	// 10Hz = 10 applied in 1000ms or every 100ms = 1000/10Hz
	// This means a stimulus is generated every 100ms which also means the
	// Inter-spike-interval (ISI) is fixed at 100ms
	m.props.Put("Hertz", 20.0)

	// Streams
	// Firing rate = spikes over an interval of time.
	m.props.Put("Firing_Rate", 0.0) // typically 0.01

	// Patterns
	m.props.Put("Poisson_Pattern_max", 0.0)
	m.props.Put("Poisson_Pattern_spread", 0.0)
	m.props.Put("Poisson_Pattern_min", 0.0)

	// Graph window properties
	// These are range values for ALL samples and are used only if RangeSync=1.0
	m.props.Put("Range_Start", 0.0)
	// End is typically set to the sim duration.
	m.props.Put("Range_End", 0.0)

	// Which synapse to focus on visually.
	m.props.Put("Active_Synapse", 0.0)
	m.props.Put("Synapse_Count", 0.0)

	// Synapse specific properties (for all synapses)
	m.props.Put("amb", 0.0) // 5
	m.props.Put("ama", 0.0) // 29

	m.props.Put("mu", 0.0)
	m.props.Put("lambda", 0.0)
	m.props.Put("alpha", 0.0)
	m.props.Put("learningRateSlow", 0.0)
	m.props.Put("learningRateFast", 0.0)

	m.props.Put("taoP", 0.0)
	m.props.Put("taoN", 0.0)
	m.props.Put("taoJ", 0.0)
	m.props.Put("distance", 0.0)

	// --------------------------------------------------------------
	// Neuron specific properties
	m.props.Put("threshold", 0.0)
	m.props.Put("RefractoryPeriod", 0.0)

	m.props.Put("nSurge", 0.0)
	m.props.Put("ntsw", 0.0)
	m.props.Put("ntao", 0.0)
	m.props.Put("ntaoS", 0.0)
	m.props.Put("ntaoJ", 0.0)
	m.props.Put("weightMin", 0.0)
	m.props.Put("weightMax", 0.0)
	m.props.Put("APMax", 0.0)

	// --------------------------------------------------------------
	// Dendrite specific properties
	m.props.Put("length", 0.0)
	m.props.Put("taoEff", 0.0)

	// ###############################################################
	// END
	// ###############################################################

	// --------------------------------------------------------------
	// Experimental properties
	m.props.Put("ExpoFunc_A", 100.0)
	m.props.Put("ExpoFunc_Tau", 100.0)
	m.props.Put("ExpoFunc_M", 5.0)
	m.props.Put("ExpoFunc_WMax", 20.0)

}

func (m *Model) Listen(msg *comm.MessageEvent) {
	// fmt.Printf("Model.Listen: %s\n", msg)

	switch msg.Target {
	case "Model":
		switch msg.Action {
		case "Set":
			_, err := strconv.ParseFloat(msg.Value, 64)
			if err == nil {
				m.SetAsFloat(msg.Field, msg.Value)
			} else {
				m.SetString(msg.Field, msg.Value)
			}
			// Notify all listeners that this property changed.
			comm.MsgBus.Send3("Model", "Data", "Changed", msg.Message, msg.ID, msg.Field, msg.Value)
			break
		case "Toggle":
			fValue := m.GetFloat(msg.Field)
			if fValue == 1.0 {
				fValue = 0.0
			} else {
				fValue = 1.0
			}
			m.SetFloat(msg.Field, fValue)
			comm.MsgBus.Send3("Model", "Data", "Changed", msg.Message, msg.ID, msg.Field, fmt.Sprintf("%0f", fValue))
			break
		}
		break
	}
}

// ----------------------------------------------------------
// Direct access
// ----------------------------------------------------------

func (m *Model) SetAsFloat(key, value string) {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	fValue, _ := strconv.ParseFloat(value, 64)
	m.props.Put(key, fValue)
}

func (m *Model) SetFloat(key string, value float64) {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	m.props.Put(key, value)
}

func (m *Model) SetString(key, value string) {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	m.props.Put(key, value)
}

func (m *Model) GetFloat(key string) float64 {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	value, _ := m.props.Get(key)
	if value == nil {
		return 0.0
	}
	return value.(float64)
}

func (m *Model) GetInt(key string) int {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	value, _ := m.props.Get(key)
	return value.(int)
}

func (m *Model) GetString(key string) string {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	value, found := m.props.Get(key)
	if found {
		return value.(string)
	}

	return ""
}

func (m *Model) GetFloatAsString(key string) string {
	m.mapMutex.Lock()
	defer m.mapMutex.Unlock()
	value, found := m.props.Get(key)
	if found {
		return fmt.Sprintf("%0.3f", value)
	}
	return ""
}
