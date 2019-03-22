package graphs

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/wdevore/Deuron5/deuron/app/events"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/gui"
	"github.com/wdevore/Deuron5/simulation/samples"
)

// SpikesAccessor provides access to series data
// type SpikesAccessor func() (x, y float64, c color.Color, done bool)

// StimulusGraph renders spikes as pixels.
// Each horizontals line is a timeline of spikes for a single
// synapse.
type StimulusGraph struct {
	BaseGraph
	gui.BaseWidget

	stimScanIdx   int
	stimScanStart int
	stimScanEnd   int

	// Spike accessor state variables
	laneOffset  float64
	colr        color.RGBA
	state       int
	poisLaneY   float64
	poisLane    *samples.SamplesLane
	poisIt      sll.Iterator
	poisScanIdx int

	stimLaneY float64
	stimLane  *samples.SamplesLane
	stimIt    sll.Iterator

	stimSamplesVisible bool
	poiSamplesVisible  bool

	noiseColor    color.RGBA
	stimulusColor color.RGBA
	unknownColor  color.RGBA

	// Synapse accessor state vars
	exciteColor  color.RGBA
	inhibitColor color.RGBA
}

// NewStimulusGraph renders spikes
func NewStimulusGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) IGraph {
	g := new(StimulusGraph)
	g.BaseWidget.Initialize(g, width, height)
	g.BaseGraph.Initialize(g.Rect, g.DC)
	g.SetGraphics(renderer, texture)

	g.stimSamplesVisible = true
	g.poiSamplesVisible = true

	g.laneOffset = 2.0
	g.state = 0
	g.exciteColor = color.RGBA{127, 127, 127, 255}
	g.inhibitColor = color.RGBA{127, 127, 255, 255}

	g.unknownColor = color.RGBA{255, 0, 0, 255}
	g.noiseColor = color.RGBA{255, 127, 0, 255}
	g.stimulusColor = color.RGBA{127, 255, 127, 255}

	return g
}

func (g *StimulusGraph) Listen(msg *comm.MessageEvent) {

	switch msg.Target {
	case "Graph":
		switch msg.Action {
		case "Toggle":
			if msg.Message == "SpikeGraph" {
				switch msg.Value {
				case "StimSamples":
					g.stimSamplesVisible = !g.stimSamplesVisible
					break
				case "PoissonSamples":
					g.poiSamplesVisible = !g.poiSamplesVisible
					break
				}
			}
			break
		}
		break
	}

	if !g.selected {
		return
	}

	poiSamples := samples.Sim.PoiSamples
	stimSamples := samples.Sim.StimSamples

	// Listen for scroll message from keymaps.go
	switch msg.Target {
	case "Key":
		switch msg.Action {
		case "Left":
			switch msg.Message {
			case "ScrollLeft":
				// Scroll window left until reach min 0. The window width
				// doesn't change.
				g.scrollLeft("SpikeGraph", poiSamples)
				g.scrollLeft("SpikeGraph", stimSamples)
				px, py := g.Position()
				g.MapToGraphSpace(g.Rect, px, py)
				break
			}
			break
		case "Right":
			switch msg.Message {
			case "ScrollRight":
				g.scrollRight("SpikeGraph", poiSamples)
				g.scrollRight("SpikeGraph", stimSamples)
				px, py := g.Position()
				g.MapToGraphSpace(g.Rect, px, py)
				break
			}
			break
		}
		break
	case "Data":
		switch msg.Action {
		case "Changed":
			switch msg.Field {
			case "Lane_Start":
				// Update sample property
				start, err := strconv.ParseFloat(msg.Value, 64)
				if err != nil {
					fmt.Printf("Error Start: %v\n", err)
				}
				poiSamples.SetRangeStart(int(start))
				stimSamples.SetRangeStart(int(start))
				break
			case "Lane_End":
				// Update sample property
				end, err := strconv.ParseFloat(msg.Value, 64)
				if err != nil {
					fmt.Printf("Error Start: %v\n", err)
				}
				poiSamples.SetRangeEnd(int(end))
				stimSamples.SetRangeEnd(int(end))
				break
			}
			break
		}
		break
	}
}

func (g *StimulusGraph) Handle(vx, vy int32, eventType events.MouseEventType) (handled bool, id int) {
	inside := gui.PointInside(vx, vy, g.Rect.X, g.Rect.Y, g.Rect.W, g.Rect.H)

	switch eventType {
	case events.MouseButton:
		if inside {
			g.selected = !g.selected
			if g.selected {
				// Send message to app.go
				comm.MsgBus.Send2("SpikeGraph", "Graph", "Selected", "Surge", fmt.Sprintf("%d", g.ID()), "")
				// Update gui Range fields the sample's range.
				samples := samples.Sim.PoiSamples
				start, end := samples.GetRange()

				// Send message to panel including the panel's id.
				comm.MsgBus.Send3("SpikeGraph", "Model", "Set", "", "", "Lane_Start", fmt.Sprintf("%d", start))
				comm.MsgBus.Send3("SpikeGraph", "Model", "Set", "", "", "Lane_End", fmt.Sprintf("%d", end))
			} else {
				comm.MsgBus.Send2("SpikeGraph", "Graph", "UnSelected", "Surge", fmt.Sprintf("%d", g.ID()), "")
			}
			return true, g.ID()
		}
		break
	case events.MouseMotion:
		handled = g.handleMotion(vx, vy, g.BaseWidget, inside)
		if handled {
			return true, g.ID()
		}
		break
	}

	return false, -1
}

// SetSeries set the chart data
func (g *StimulusGraph) SetSeries(accessor SeriesAccessor) {
}

// Destroy release resources
func (g *StimulusGraph) Destroy() {
}

// DrawAt renders graph to texture
func (g *StimulusGraph) Draw() {
	g.BaseGraph.Draw(g.Rect, g.DC)

	// -------------------------------------------
	// Draw data
	// -------------------------------------------

	// TODO Draw colored horizontal bars based on the synapse type.
	winX := float64(0)

	// Draw noise spikes.
	if g.poiSamplesVisible {
		px, py, c, spiked := g.poissonAccessor()

		// Map from sample-space to unit-space first
		uspX := g.Linear(float64(g.scanStart), float64(g.scanEnd), px)

		// Map from unit-space to window-space
		// 0 is the min because we have already translated the origin above.
		winX = g.Lerp(0, g.upperX, uspX)

		for spiked > 0 {
			if spiked == 1 {
				g.DC.SetColor(c)
				// g.DC.DrawPoint(px, py, 1.0)
				tx, ty := g.DC.TransformPoint(winX, py)
				g.DC.SetPixel(int(tx), int(ty))
				g.DC.Fill()
			}

			px, py, c, spiked = g.poissonAccessor()
			uspX = g.Linear(float64(g.scanStart), float64(g.scanEnd), px)
			winX = g.Lerp(0, g.upperX, uspX)
		}
	}

	// Draw stimulus spikes.
	if g.stimSamplesVisible {
		sx, sy, sc, sSpiked := g.stimAccessor()
		uspX := g.Linear(float64(g.stimScanStart), float64(g.stimScanEnd), sx)
		winX = g.Lerp(0, g.upperX, uspX)

		for sSpiked > 0 {
			if sSpiked == 1 {
				g.DC.SetColor(sc)
				tx, ty := g.DC.TransformPoint(winX, sy)
				g.DC.SetPixel(int(tx), int(ty))
				g.DC.Fill()
			}

			sx, sy, sc, sSpiked = g.stimAccessor()
			uspX = g.Linear(float64(g.stimScanStart), float64(g.stimScanEnd), sx)
			winX = g.Lerp(0, g.upperX, uspX)
		}
	}

	// -------------------------------------------
	// Draw labels and ruler marks
	// -------------------------------------------
	g.drawTitle(g.Rect, g.DC)

	if g.poiSamplesVisible || g.stimSamplesVisible {
		winX = g.drawVerticalTimeBar(g.Rect, g.DC)
		g.drawMouseInfo(winX, g.Rect, g.DC)
	}

	g.postDraw(g.Rect)
}

func (g *StimulusGraph) Check() bool {
	if samples.Sim.PoiSamples == nil {
		return false
	}

	poiLanes := samples.Sim.PoiSamples.GetLanes()
	g.poisIt = poiLanes.Iterator()

	if g.poisIt.First() {
		g.poisLane = g.poisIt.Value().(*samples.SamplesLane)
		g.setScanWindow(samples.Sim.PoiSamples)

		stimLanes := samples.Sim.StimSamples.GetLanes()
		g.stimIt = stimLanes.Iterator()

		if g.stimIt.First() {
			g.stimLane = g.stimIt.Value().(*samples.SamplesLane)
			fmt.Printf("%d\n", g.stimLane.Id)

			g.setStimScanWindow(samples.Sim.StimSamples)
			return true
		}

		return false
	}

	return false
}

func (g *StimulusGraph) setScanWindow(lanes *samples.Samples) {
	// If RangeSync is enabled then we use the Model's range
	sync := deuron.SimModel.GetFloat("RangeSync")

	if sync == 1 {
		g.scanStart = int(deuron.SimModel.GetFloat("Range_Start"))
		g.scanEnd = int(deuron.SimModel.GetFloat("Range_End"))
	} else {
		g.scanStart, g.scanEnd = lanes.GetRange()
	}

	g.scanIdx = g.scanStart
}

func (g *StimulusGraph) setStimScanWindow(lanes *samples.Samples) {
	// If RangeSync is enabled then we use the Model's range
	sync := deuron.SimModel.GetFloat("RangeSync")

	if sync == 1 {
		g.stimScanStart = int(deuron.SimModel.GetFloat("Range_Start"))
		g.stimScanEnd = int(deuron.SimModel.GetFloat("Range_End"))
	} else {
		g.stimScanStart, g.stimScanEnd = lanes.GetRange()
	}

	g.stimScanIdx = g.stimScanStart
}

func (g *StimulusGraph) poissonAccessor() (x, y float64, c color.Color, spiked int) {
	if g.scanIdx >= g.scanEnd {
		g.setScanWindow(samples.Sim.PoiSamples)
		if !g.poisIt.Next() {
			g.poisLaneY = 0
			return 0, 0, nil, 0
		}

		g.poisLaneY += g.laneOffset + 1
		g.poisLane = g.poisIt.Value().(*samples.SamplesLane)
	}

	spike := g.poisLane.Values[g.scanIdx]

	if spike.Value == 1 {
		g.state = 1
	} else {
		g.state = 2
	}

	g.scanIdx++
	return spike.Time, g.poisLaneY, g.noiseColor, g.state
}

func (g *StimulusGraph) stimAccessor() (x, y float64, c color.Color, spiked int) {
	if g.stimScanIdx >= g.stimScanEnd {
		g.setStimScanWindow(samples.Sim.StimSamples)
		if !g.stimIt.Next() {
			g.stimLaneY = 0
			return 0, 0, nil, 0
		}

		g.stimLaneY += g.laneOffset + 1
		g.stimLane = g.stimIt.Value().(*samples.SamplesLane)
		fmt.Printf("%d\n", g.stimLane.Id)
	}

	spike := g.stimLane.Values[g.stimScanIdx]
	// fmt.Printf("s: %v\n", spike)

	if spike.Value == 1 {
		g.state = 1
	} else {
		g.state = 2
	}

	g.stimScanIdx++
	return spike.Time, g.stimLaneY, g.stimulusColor, g.state
}
