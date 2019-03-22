package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
	"github.com/wdevore/Deuron5/deuron/app/graphs"
	"github.com/wdevore/Deuron5/deuron/app/gui"
	"github.com/wdevore/Deuron5/simulation/runreset"
	"github.com/wdevore/Deuron5/simulation/samples"
)

const (
	width  = 1000 * 2
	height = 800 * 2
	xpos   = 100
	ypos   = 0
	FPS    = 30 // 60 or 30 fps
)

// App shows the plots and graphs.
// It receives commands for graphing and viewing various graphs.
type App struct {
	window   *sdl.Window
	surface  *sdl.Surface
	renderer *sdl.Renderer

	// The graphs and text are rendered to this texture.
	texture *sdl.Texture

	running      bool
	shuttingDown bool

	opened  bool
	created bool

	nFont *Font

	dynaText     *DynaText
	txtSimStatus *Field
	txtTime      *Field

	// Fields. Only one is displayed at any one time.
	// When a field changes it transmits this to the simulation.
	txtActiveProperty *Field

	graphs *sll.List

	dirty bool

	// expoGraph   graphs.IGraph

	//keymapBar *KeymapBar
	gui *gui.Gui

	simType string // "runreset" or "continous"
	target  string // Path to simulation

	// comm channel to simulation
	statusComm chan string

	// simulation = run_reset.go
	simulation deuron.ISimulation

	// Commands
	thing    string
	property string

	// 10 keyboard map layouts
	mode   string // "Main", "Entry"
	mapIdx int

	// -----------------------------------------------------------------
	// Keymap fields
	// -----------------------------------------------------------------
	// Current mouse position
	mouseX, mouseY  int32
	controlID       string
	controlMsg      string
	controlOrgValue string
	controlField    string
	currentValue    string

	jsonMap map[string]interface{}
}

// NewApp creates a new App and initializes it.
func NewApp() *App {
	ap := new(App)
	ap.opened = false
	ap.simType = "runreset"
	ap.mode = "Main"
	ap.dirty = true
	return ap
}

// Open shows the App and begin event polling
// (host deuron.IHost)
func (ap *App) Open() {
	ap.Load("neuron.json")

	ap.initialize()

	ap.shuttingDown = false

	ap.opened = true
}

func (ap *App) subscribeToMsgBus() {
	comm.MsgBus.Subscribe(ap)

	ap.gui.SubscribeToMsgBus()

	comm.MsgBus.Subscribe(deuron.SimModel)

	it := ap.graphs.Iterator()
	for it.Next() {
		listener := it.Value().(comm.IMessageListener)
		comm.MsgBus.Subscribe(listener)
	}

}

func (ap *App) Handle(vx, vy int32, eventType events.MouseEventType) (handled bool) {
	handled = ap.gui.Handle(vx, vy, eventType)
	if handled {
		return true
	}

	it := ap.graphs.Iterator()
	for it.Next() {
		widget := it.Value().(gui.IWidget)
		handled, _ = widget.Handle(vx, vy, eventType)
		if handled {
			return true
		}
	}

	return false
}

func (ap *App) Listen(msg *comm.MessageEvent) {
	// fmt.Printf("App.Listen: %s\n", msg)

	switch msg.Source {
	case "Model", "Gui":
		switch msg.Target {
		case "Data", "Button", "ToggleButton":
			switch msg.Action {
			case "Changed":
				ap.dirty = true
				break
			}
			break
		}
		break
	}

	switch msg.Target {
	case "Model":
		switch msg.Action {
		case "Change":
			switch msg.Message {
			case "Status":
				ap.txtSimStatus.SetValue(msg.Value)
				return
			}
			break
		case "Data":
			switch msg.Message {
			case "Changed":
				ap.dirty = true
				return
			}
			break
		}
		break
	case "App":
		switch msg.Action {
		case "Edit":
			switch msg.Message {
			case "Start":
				ap.controlID = msg.ID
				ap.controlField = msg.Field
				ap.controlMsg = msg.Message
				ap.controlOrgValue = msg.Value
				ap.currentValue = ""
				ap.dirty = true
				return
			}
			break
		case "Command":
			switch msg.Message {
			case "RunPause":
				// fmt.Println("---|||| Running and pausing simulation ||||---")
				// Run a single pass through a complete simulation then pause.
				ap.Create()
				ap.RunPause()
				ap.dirty = true
				return
			case "Save":
				ap.Save("neuron.json")
				return
			case "Load":
				ap.Load("neuron.json")
				return
			}
			break
		}
		break
	case "Graph":
		switch msg.Action {
		case "Selected":
			// Unselect other graphs
			ap.controlID = msg.ID
			ap.currentValue = ""
			ap.dirty = true
			return
		case "UnSelected":
			ap.controlID = ""
			ap.dirty = true
			return
		}
		break
	}

	if ap.simulation != nil {
		ap.simulation.SendEvent(msg)
	}
}

// SetFont sets the font based on path and size.
func (ap *App) SetFont(fontPath string, size int) {
	ap.nFont = NewFont(fontPath, size)
}

func (ap *App) SetText(field, value string) {
	ap.txtActiveProperty.SetName(field + ": ")
	ap.txtActiveProperty.SetValue(value)
}

func (ap *App) SetValue(value string) {
	ap.txtActiveProperty.SetValue(value)
}

// Run starts the polling event loop. This must execute on the main thread.
func (ap *App) Run() {
	ap.running = true

	for ap.running {
		sdl.PumpEvents()

		if ap.dirty {
			ap.Update()
			ap.dirty = false
		}

		// sdl.Delay(17)
		time.Sleep(time.Millisecond * (1000 / FPS))
	}

	ap.shutdown()
}

func (ap *App) Update() {
	if ap.shuttingDown {
		return
	}

	ap.renderer.Clear()

	ap.updateGraphs()

	ap.txtSimStatus.Draw()
	ap.txtTime.Draw()
	ap.txtActiveProperty.Draw()

	// Draw keymap overlays
	//ap.keymapBar.DrawAt(500, 200)
	ap.gui.Draw()

	ap.renderer.Present()
	// ap.window.UpdateSurface()
}

// Quit stops the gui from running, effectively shutting it down.
func (ap *App) Quit() {
	ap.running = false
}

// Close closes the App.
// Be sure to setup a "defer x.Close()"
func (ap *App) Close() {
	if !ap.opened {
		return
	}

	log.Println("\nClosing App...")
	comm.MsgBus.Close()

	// log.Println("Destroying text(s)")
	ap.txtSimStatus.Destroy()
	ap.txtActiveProperty.Destroy()
	ap.txtTime.Destroy()

	// log.Println("Destroying font")
	ap.nFont.Destroy()

	// log.Println("Destroying graphs")

	it := ap.graphs.Iterator()
	for it.Next() {
		graph := it.Value().(graphs.IGraph)
		graph.Destroy()
	}

	// log.Println("Destroying texture")
	err := ap.texture.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Destroying renderer")
	err = ap.renderer.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	// log.Println("Shutting down App")
	err = ap.window.Destroy()
	sdl.Quit()

	if err != nil {
		log.Fatal(err)
	}
}

func (ap *App) initialize() {
	var err error

	err = sdl.Init(sdl.INIT_TIMER | sdl.INIT_VIDEO | sdl.INIT_EVENTS)
	if err != nil {
		panic(err)
	}

	ap.window, err = sdl.CreateWindow("Deuron5 Graph", xpos, ypos, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	ap.surface, err = ap.window.GetSurface()
	if err != nil {
		panic(err)
	}

	// ap.renderer, err = sdl.CreateSoftwareRenderer(ap.surface)
	ap.renderer, err = sdl.CreateRenderer(ap.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	ap.texture, err = ap.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}
}

// Configure view with draw objects
func (ap *App) Configure() {
	fmt.Println("App configuring...")

	ap.createGraphs()
	ap.createGui()

	sdl.SetEventFilterFunc(ap.filterEvent, nil)

	// sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")
	ap.renderer.SetDrawColor(64, 64, 64, 255)

	ap.dynaText = NewDynaText(ap.nFont, ap.renderer)

	ap.txtSimStatus = NewField(ap.nFont, ap.renderer, ap.dynaText)
	ap.txtSimStatus.SetNameColor(sdl.Color{R: 127, G: 64, B: 0, A: 255})
	ap.txtSimStatus.SetName("Status: ")
	ap.txtSimStatus.SetValue(deuron.SimModel.GetString("Status"))

	ap.txtTime = NewField(ap.nFont, ap.renderer, ap.dynaText)
	ap.txtTime.SetName("")
	ap.txtTime.SetValue("0")
	ap.txtTime.SetPosition(200, 0)

	ap.txtActiveProperty = NewField(ap.nFont, ap.renderer, ap.dynaText)
	ap.txtActiveProperty.SetPosition(5, 30)

	ap.subscribeToMsgBus()
}

func (ap *App) createGui() {
	ap.gui = gui.NewGui(ap.renderer, ap.texture, 1000, 50)
}

// --------------------------------------------------------------------
// Graphs
// --------------------------------------------------------------------

func (ap *App) createGraphs() {
	ap.graphs = sll.New()

	graphIDs := 10000
	graph := graphs.NewStimulusScatterGraph(ap.renderer, ap.texture, 2000, 100)
	graph.SetName("Stimulus")
	widget := graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, 75)
	ap.graphs.Add(graph)

	_, y := widget.Position()
	graphIDs++
	graph = graphs.NewSurgeGraph(ap.renderer, ap.texture, 2000, 200)
	graph.SetName("Synapse Surge")
	widget = graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, y+100)
	ap.graphs.Add(graph)

	_, y = widget.Position()
	graphIDs++
	graph = graphs.NewPspGraph(ap.renderer, ap.texture, 2000, 200)
	graph.SetName("Synapse PSP")
	widget = graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, y+200)
	ap.graphs.Add(graph)

	_, y = widget.Position()
	graphIDs++
	graph = graphs.NewWeightGraph(ap.renderer, ap.texture, 2000, 200)
	graph.SetName("Synapse Weight")
	widget = graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, y+200)
	ap.graphs.Add(graph)

	_, y = widget.Position()
	graphIDs++
	graph = graphs.NewNeuronPspGraph(ap.renderer, ap.texture, 2000, 200)
	graph.SetName("Neuron PSP")
	widget = graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, y+200)
	ap.graphs.Add(graph)

	_, y = widget.Position()
	graphIDs++
	graph = graphs.NewPostSpikeGraph(ap.renderer, ap.texture, 2000, 50)
	graph.SetName("Post Spike")
	widget = graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, y+200)
	ap.graphs.Add(graph)

	_, y = widget.Position()
	graphIDs++
	graph = graphs.NewNeuronAPGraph(ap.renderer, ap.texture, 2000, 100)
	graph.SetName("Neuron AP fast")
	widget = graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, y+50)
	ap.graphs.Add(graph)

	_, y = widget.Position()
	graphIDs++
	graph = graphs.NewNeuronAPSlowGraph(ap.renderer, ap.texture, 2000, 100)
	graph.SetName("Neuron AP slow")
	widget = graph.(gui.IWidget)
	widget.SetID(graphIDs)
	widget.SetPos(0, y+100)
	ap.graphs.Add(graph)

	// _, y = widget.Position()
	// graphIDs++
	// graph = graphs.NewDTGraph(ap.renderer, ap.texture, 2000, 50)
	// graph.SetName("dt")
	// widget = graph.(gui.IWidget)
	// widget.SetID(graphIDs)
	// widget.SetPos(0, y+200)
	// ap.graphs.Add(graph)

	// _, y = widget.Position()
	// graphIDs++
	// graph = graphs.NewNeuronDtGraph(ap.renderer, ap.texture, 2000, 50)
	// graph.SetName("Neuron dt")
	// widget = graph.(gui.IWidget)
	// widget.SetID(graphIDs)
	// widget.SetPos(0, y+50)
	// ap.graphs.Add(graph)

	// ap.expoGraph = graphs.NewExpoGraph(ap.renderer, ap.texture, 1024, 600)
}

func (ap *App) updateGraphs() {
	if samples.Sim != nil {
		it := ap.graphs.Iterator()
		for it.Next() {
			graph := it.Value().(graphs.IGraph)
			draw := graph.Check()
			if draw {
				graph.Draw()
			}
		}
	}
}

func (ap *App) Load(file string) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Loading model from (%s)\n", file)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	ap.jsonMap = make(map[string]interface{})

	err = json.Unmarshal(byteValue, &ap.jsonMap)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Directly load model rather than send messages to listeners
	// who don't exist yet.
	m := deuron.SimModel

	// Calculate sim-duration (aka sample size) based on TimeStep and duration
	duration := ap.jsonMap["Duration"].(float64)
	timeStep := ap.jsonMap["TimeStep"].(float64)
	simDuration := duration * 1000000.0 / timeStep
	fmt.Printf("Sample size: %d\n", int(simDuration))

	m.SetFloat("TimeStep", timeStep) // microseconds
	m.SetFloat("Duration", duration) // seconds
	m.SetFloat("Samples", simDuration)

	m.SetFloat("Range_Start", ap.jsonMap["RangeStart"].(float64))
	m.SetFloat("Range_End", ap.jsonMap["RangeEnd"].(float64))
	m.SetString("Stimulus", ap.jsonMap["Stimulus"].(string))
	m.SetFloat("Synapse_Count", ap.jsonMap["Synapse_Count"].(float64))

	m.SetFloat("weightMin", ap.jsonMap["weightMin"].(float64))
	m.SetFloat("weightMax", ap.jsonMap["weightMax"].(float64))

	fmt.Println("Loaded")
}

func (ap *App) Save(file string) {
	// x := map[string]interface{}{"a": 1.1}
	// y := []interface{}{1, 2, 3}
	// ap.jsonMap["Help"] = y
	if ap.simulation == nil {
		fmt.Println("Missing simulation object.")
		return
	}

	ap.saveSimSettings(file)

	ap.saveSimulation()

	fmt.Println("---- Saved ----.")
}

func (ap *App) saveSimSettings(file string) {
	mo := deuron.SimModel

	m := map[string]interface{}{
		"weightMin":     mo.GetFloat("weightMin"),
		"weightMax":     mo.GetFloat("weightMax"),
		"RangeStart":    mo.GetFloat("Range_Start"),
		"RangeEnd":      mo.GetFloat("Range_End"),
		"Duration":      mo.GetFloat("Duration"),
		"TimeStep":      mo.GetFloat("TimeStep"),
		"Synapse_Count": mo.GetFloat("Synapse_Count"),
		"Stimulus":      mo.GetString("Stimulus"),
	}

	jsonString, err := json.MarshalIndent(m, "", "  ")

	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("json: \n%s\n", jsonString)

	fmt.Printf("Writing simulation settings to (%s)\n", file)

	err = ioutil.WriteFile(file, jsonString, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func (ap *App) saveSimulation() {
	// Get json map
	nSim := ap.simulation.ToJSON()

	// Serialize
	jsonString, err := json.MarshalIndent(nSim, "", "  ")

	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Printf("json: \n%s\n", jsonString)
	file := deuron.SimModel.GetString("Stimulus")

	fileName := "./stimulus/" + file + ".json"

	fmt.Printf("Writing simulation to (%s)\n", fileName)

	err = ioutil.WriteFile(fileName, jsonString, 0644)

	if err != nil {
		fmt.Println(err)
		return
	}
}

// ----------------------------------------------------------------
// Commands
// ----------------------------------------------------------------

func (ap *App) Go() {
	go ap.doit()
}

func (ap *App) Create() {
	if ap.created {
		return
	}

	// Connect
	ap.connect()

	ap.simulation.Create()

	ap.gui.Refresh()

	ap.created = true

	go ap.pollForMessage()
}

func (ap *App) Step() {
	ap.simulation.Step()
}

func (ap *App) RunPause() {
	// Once AutoSync is complete we can run in a "go" routine.
	ap.simulation.RunPause()
}

func (ap *App) Reset() {
	ap.simulation.Reset()
}

func (ap *App) Pause() {
}

// Command handles messages from the console.
func (ap *App) Command(args []string) {
	switch args[0] {
	case "quit":
		ap.shuttingDown = true
		ap.shutdown()
	case "set":
		ap.target = args[1]
	case "type":
		ap.simType = args[1]
		fmt.Printf("Type switched to `%s`\n", ap.simType)
	case "con":
		ap.connect()
	case "ping":
		ap.simulation.Send("ping")
	case "start":
		ap.start()
	case "stop":
		if ap.simulation == nil {
			panic("Not connected. Please connect first. Use 'help'.")
		}
		ap.simulation.Send("stop")
		response := <-ap.statusComm
		if response == "Stopped" {
			fmt.Println("Sim requested to stop")
		}
		break
	case "prop":
		// A command relating to a property
		ap.simulation.Command(args)
		break
	}
}

func (ap *App) doit() {
	// Connect
	ap.connect()

	// Load
	ap.simulation.Send("load")
	response := <-ap.statusComm // wait for response

	if response != "loaded" {
		panic("Unable to load parameters")
	}

	fmt.Println("Loaded")

	deuron.SimModel.SetString("Status", "Starting...")
	ap.txtSimStatus.SetValue(deuron.SimModel.GetString("Status"))

	ap.start()
}

func (ap *App) start() {
	// We start async because we can't lock the app thread
	// from receiveing system events (ex: keyboard)
	go ap.pollForMessage()

	// This will cause the sim to start the simulation in a coroutine.
	duration := deuron.SimModel.GetFloatAsString("Samples")

	ap.simulation.Send(fmt.Sprintf("start %s", duration))
}

// Runs in a coroutine.
func (ap *App) pollForMessage() {
	var response string
	fmt.Print("Waiting for messages from simulation coroutine...\n")

	for response != "Stopped" {
		response = <-ap.statusComm // Wait for response

		if response == "GuiRefesh" {
			// fmt.Println("Refreshing gui")
			ap.gui.Refresh()
		}
		// deuron.SimModel.SetString("Status", response)
		ap.txtTime.SetValue(response)
	}

	fmt.Printf("Polling exited from: (%s)\n", response)
}

func (ap *App) connect() {
	fmt.Printf("Connecting to `%s`...\n", ap.simType)

	if ap.simulation == nil {
		fmt.Println("Creating sim")
		ap.simulation = runreset.NewRunResetSim()

		fmt.Println("Creating comm channels")

		ap.statusComm = make(chan string)
		ap.simulation.Connect(ap.statusComm)

		// Test connection to sim
		ap.simulation.Send("ping")
		response := <-ap.statusComm
		if response == "pong" {
			fmt.Println("Connected.")
		} else {
			fmt.Printf("Sim didn't respond to connection correctly (%s)\n", response)
		}
	}
}

func (ap *App) shutdown() {
	fmt.Println("Shutting down...")
	ap.shuttingDown = true

	if ap.simulation != nil {
		fmt.Println("Sending simulation the `stop` command...")
		ap.simulation.Send("stop")
	}

	ap.Quit()

	fmt.Println("Done.")
}
