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

type PostSpikeGraph struct {
	BaseGraph
	gui.BaseWidget

	graphIt sll.Iterator

	// Synapse accessor state vars
	lineColor    color.RGBA
	inhibitColor color.RGBA
}

// NewPostSpikeGraph renders spikes
func NewPostSpikeGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) IGraph {
	g := new(PostSpikeGraph)
	g.BaseWidget.Initialize(g, width, height)
	g.BaseGraph.Initialize(g.Rect, g.DC)
	g.SetGraphics(renderer, texture)

	g.lineColor = color.RGBA{127, 255, 255, 255}

	return g
}

func (g *PostSpikeGraph) Listen(msg *comm.MessageEvent) {

	if !g.selected {
		return
	}

	samples := samples.Sim.CellSamples

	// Listen for scroll message from keymaps.go
	switch msg.Target {
	case "Key":
		switch msg.Action {
		case "Left":
			switch msg.Message {
			case "ScrollLeft":
				// Scroll window left until reach min 0. The window width
				// doesn't change.
				g.scrollLeft("PostSpikeGraph", samples)
				px, py := g.Position()
				g.MapToGraphSpace(g.Rect, px, py)
				break
			}
			break
		case "Right":
			switch msg.Message {
			case "ScrollRight":
				g.scrollRight("PostSpikeGraph", samples)
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
				samples.SetRangeStart(int(start))
				break
			case "Lane_End":
				// Update sample property
				end, err := strconv.ParseFloat(msg.Value, 64)
				if err != nil {
					fmt.Printf("Error Start: %v\n", err)
				}
				samples.SetRangeEnd(int(end))
				break
			}
			break
		}
		break
	}
}

func (g *PostSpikeGraph) Handle(vx, vy int32, eventType events.MouseEventType) (handled bool, id int) {
	inside := gui.PointInside(vx, vy, g.Rect.X, g.Rect.Y, g.Rect.W, g.Rect.H)

	switch eventType {
	case events.MouseButton:
		if inside {
			g.selected = !g.selected
			if g.selected {
				// Send message to app.go
				comm.MsgBus.Send2("PostSpikeGraph", "Graph", "Selected", "Surge", fmt.Sprintf("%d", g.ID()), "")
				// Update gui Range fields the sample's range.
				samples := samples.Sim.CellSamples
				start, end := samples.GetRange()

				// Send message to panel including the panel's id.
				comm.MsgBus.Send3("PostSpikeGraph", "Model", "Set", "", "", "Lane_Start", fmt.Sprintf("%d", start))
				comm.MsgBus.Send3("PostSpikeGraph", "Model", "Set", "", "", "Lane_End", fmt.Sprintf("%d", end))
			} else {
				comm.MsgBus.Send2("PostSpikeGraph", "Graph", "UnSelected", "Surge", fmt.Sprintf("%d", g.ID()), "")
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
func (g *PostSpikeGraph) SetSeries(accessor SeriesAccessor) {
}

// Destroy release resources
func (g *PostSpikeGraph) Destroy() {
}

// DrawAt renders graph to texture
func (g *PostSpikeGraph) Draw() {
	g.BaseGraph.Draw(g.Rect, g.DC)

	// -------------------------------------------
	// Draw data
	// -------------------------------------------

	// TODO Draw colored horizontal bars based on the synapse type.
	winX := float64(0)
	winY := float64(0)

	px, py, c, more := g.dataAccessor()

	// Map from sample-space to unit-space first
	uspX := g.Linear(float64(g.scanStart), float64(g.scanEnd), px)
	// uspY := g.Linear(0.0, 1.0, 1.0)

	// Map from unit-space to window-space
	// 0 is the min because we have already translated the origin above.
	winX = g.Lerp(0, g.upperX, uspX)
	winY = g.Lerp(0, g.upperY, 1.0)

	for more > 0 {
		if py > 0 {
			g.DC.SetColor(c)
			g.DC.MoveTo(winX, 0)
			g.DC.LineTo(winX, winY)
			g.DC.Stroke()
		}

		px, py, c, more = g.dataAccessor()
		uspX = g.Linear(float64(g.scanStart), float64(g.scanEnd), px)
		winX = g.Lerp(0, g.upperX, uspX)
	}

	// -------------------------------------------
	// Draw labels and ruler marks
	// -------------------------------------------
	g.drawTitle(g.Rect, g.DC)

	winX = g.drawVerticalTimeBar(g.Rect, g.DC)
	g.drawMouseInfo(winX, g.Rect, g.DC)

	g.postDraw(g.Rect)
}

func (g *PostSpikeGraph) Check() bool {
	if samples.Sim.CellSamples == nil {
		return false
	}

	lanes := samples.Sim.CellSamples.GetLanes()

	g.graphIt = lanes.Iterator()

	if g.graphIt.First() {
		lane, _ := lanes.Get(0)
		g.activeLane = lane.(*samples.SamplesLane)
		g.setScanWindow(samples.Sim.CellSamples)
		return true
	}

	return false
}

func (g *PostSpikeGraph) setScanWindow(lanes *samples.Samples) {
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

func (g *PostSpikeGraph) dataAccessor() (x, y float64, c color.Color, spiked int) {
	// Return values until we reach the end of the current active lane.

	if g.scanIdx >= g.scanEnd {
		// No more samples in the lane.
		g.setScanWindow(samples.Sim.CellSamples)
		return 0, 0, nil, 0
	}

	sample := g.activeLane.Values[g.scanIdx]

	g.scanIdx++

	if sample.Value != nil {
		// fmt.Printf("value %0.1f\n", sample.Value.(float64))
		return sample.Time, sample.Value.(float64), g.lineColor, 1
	}

	return 0, 0, nil, 0
}
