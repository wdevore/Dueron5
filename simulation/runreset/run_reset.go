package runreset

import (
	"fmt"
	"strings"

	"github.com/wdevore/Deuron5/deuron/app/comm"

	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/simulation/samples"
)

var TimeStep = 1.0 // milliseconds

/*
This simulation simulates a single neuron by repeatedly
applying stimulus for a time course and then resetting and repeating.
This type of simulation is for tuning the neuron's parameters such that
the neuron maintains and operates within a critical functional region.
*/

type RunResetSim struct {
	statusChannel chan string

	stopped bool

	workingPath string

	// Sim ticks at TimeStep(ms) resolution
	t float64

	sim *Simulation
}

func NewRunResetSim() deuron.ISimulation {
	s := new(RunResetSim)
	s.stopped = true
	return s
}

func (s *RunResetSim) Connect(statusChannel chan string) {
	// Send message back to the App/Viewer in the
	// responseLoop coroutine.
	s.statusChannel = statusChannel
}

func (s *RunResetSim) Command(args []string) {
	// fmt.Printf("Send msg: %s\n", msg)
	switch args[0] {
	case "start":
		s.stopped = false
		s.Start()
	case "stop":
		s.stopped = true
	case "ping":
		go s.respond("pong")
	}
}

func (s *RunResetSim) Send(msg string) {
	args := strings.Split(msg, " ")
	s.Command(args)
}

// Sends a msg back through channel async
func (s *RunResetSim) respond(msg string) {
	s.statusChannel <- msg
}

func (s *RunResetSim) Start() {
	s.Create()
	fmt.Println("Starting...")

	// Start the simulation loop in a coroutine.
	go s.run()
	deuron.SimModel.SetString("Status", "Running...")
}

func (s *RunResetSim) Create() {
	fmt.Println("Creating...")
	s.t = 0.0

	s.sim = NewSimulation(s.statusChannel)

	// Setup the neuron while connecting the noise and stimulus.
	synCnt := s.sim.initialize()

	duration := deuron.SimModel.GetFloat("Samples")

	sampleSize := int(duration)

	fmt.Printf("Syn cnt: %d, duration: %d\n", synCnt, sampleSize)

	// The samples is where we collect all the data.
	samples.Sim = samples.NewSamplesCollection(synCnt, sampleSize)

	fmt.Println("Created.")
}

// This runs in a "Go"routine.
func (s *RunResetSim) run() {
	fmt.Println("RunReset: run() loop begining")
	// Run the sim for a fixed amount of time and then reset.
	duration := deuron.SimModel.GetFloat("Samples")

	for !s.stopped {
		if s.t >= duration {
			s.Reset()
		} else {
			s.Step()
		}
	}

	fmt.Println("RunReset: run() loop exited")
	s.respond("Stopped")
}

func (s *RunResetSim) Reset() {
	// Reset
	s.t = 0.0

	// Reset random seeds.
	s.sim.reset()
}

func (s *RunResetSim) Step() {
	s.sim.simulate(s.t)
	s.t += TimeStep
}

// This can run in a goroutine or not.
func (s *RunResetSim) RunPause() {
	s.Reset()
	duration := deuron.SimModel.GetFloat("Samples")

	fmt.Println("Starting run...")
	for s.t < duration {
		// Run
		s.Step()
	}
	fmt.Println("Run complete.")

	s.sim.PostProcess()
}

func (s *RunResetSim) SendEvent(event *comm.MessageEvent) {
	s.sim.SendEvent(event)
}

func (s *RunResetSim) ToJSON() interface{} {
	return s.sim.ToJSON()
}

func (s *RunResetSim) Load(json interface{}) {
	s.sim.Load(json)
}
