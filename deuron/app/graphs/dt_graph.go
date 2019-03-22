package graphs

import (
	"fmt"
	"image/color"

	"github.com/wdevore/Deuron5/deuron/app/events"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/gui"
	"github.com/wdevore/Deuron5/simulation/samples"
)

type DTGraph struct {
	BaseGraph
	gui.BaseWidget

	lineColor    color.RGBA
	altLineColor color.RGBA
	even         bool

	prevWinX, prevWinY float64

	graphIt sll.Iterator
}

func NewDTGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) IGraph {
	g := new(DTGraph)
	g.BaseWidget.Initialize(g, width, height)
	g.BaseGraph.Initialize(g.Rect, g.DC)
	g.SetGraphics(renderer, texture)

	g.lineColor = color.RGBA{255, 127, 0, 255}
	g.altLineColor = color.RGBA{255, 160, 0, 255}

	return g
}

func (g *DTGraph) Listen(msg *comm.MessageEvent) {
	if !g.selected {
		return
	}

	samples := samples.Sim.DtSamples

	px, py := g.Position()
	handled := g.handleScroll(msg, samples, px, py, g.Rect)

	if handled {
		return
	}

	handled = g.handleRange(msg, samples)

	if handled {
		return
	}
}

func (g *DTGraph) Handle(vx, vy int32, eventType events.MouseEventType) (handled bool, id int) {
	inside := gui.PointInside(vx, vy, g.Rect.X, g.Rect.Y, g.Rect.W, g.Rect.H)

	switch eventType {
	case events.MouseButton:
		if inside {
			g.selected = !g.selected
			if g.selected {
				// Send message to app.go
				comm.MsgBus.Send2("DTGraph", "Graph", "Selected", "Surge", fmt.Sprintf("%d", g.ID()), "")
				// Update gui Range fields the sample's range.
				samples := samples.Sim.DtSamples
				start, end := samples.GetRange()

				// Send message to panel including the panel's id.
				comm.MsgBus.Send3("DTGraph", "Model", "Set", "", "", "Lane_Start", fmt.Sprintf("%d", start))
				comm.MsgBus.Send3("DTGraph", "Model", "Set", "", "", "Lane_End", fmt.Sprintf("%d", end))
			} else {
				comm.MsgBus.Send2("DTGraph", "Graph", "UnSelected", "Surge", fmt.Sprintf("%d", g.ID()), "")
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

func (g *DTGraph) SetSeries(accessor SeriesAccessor) {
}

// Destroy release resources
func (g *DTGraph) Destroy() {
}

func (g *DTGraph) Prep() {
}

// DrawAt renders graph to texture
func (g *DTGraph) Draw() {
	g.BaseGraph.Draw(g.Rect, g.DC)

	// -------------------------------------------
	// Draw data
	// -------------------------------------------
	px, py, _, more := g.dataAccessor()

	// Map from sample-space to unit-space first
	uspX := g.Linear(float64(g.scanStart), float64(g.scanEnd), px)
	uspY := g.Linear(g.activeLane.Min, g.activeLane.Max, py)

	// Map from unit-space to window-space
	// 0 is the min because we have already translated the origin above.
	winX := g.Lerp(0, g.upperX, uspX)
	winY := g.Lerp(0, g.upperY, uspY)

	g.even = true

	g.DC.SetColor(g.lineColor)

	px, py, _, more = g.dataAccessor()
	for more > 0 {
		g.DC.MoveTo(winX, winY)

		if g.even {
			g.DC.SetColor(g.altLineColor)
		} else {
			g.DC.SetColor(g.lineColor)
		}

		uspX := g.Linear(float64(g.scanStart), float64(g.scanEnd), px)
		uspY = g.Linear(g.activeLane.Min, g.activeLane.Max, py)
		winX = g.Lerp(0, g.upperX, uspX)
		winY = g.Lerp(0, g.upperY, uspY)

		g.DC.LineTo(winX, winY)

		px, py, _, more = g.dataAccessor()
		g.even = !g.even
		g.DC.Stroke()
	}

	// -------------------------------------------
	// Draw labels and ruler marks
	// -------------------------------------------
	g.drawTitle(g.Rect, g.DC)
	g.drawMinMax(g.activeLane.Min, g.activeLane.Max, g.Rect, g.DC)

	winX = g.drawVerticalTimeBar(g.Rect, g.DC)
	g.drawMouseInfo(winX, g.Rect, g.DC)

	g.postDraw(g.Rect)
}

func (g *DTGraph) Check() bool {
	if samples.Sim.DtSamples == nil {
		return false
	}

	lanes := samples.Sim.DtSamples.GetLanes()

	g.graphIt = lanes.Iterator()

	if g.graphIt.First() {
		activeSynID := int(deuron.SimModel.GetFloat("Active_Synapse"))
		lane, _ := lanes.Get(activeSynID)
		g.activeLane = lane.(*samples.SamplesLane)
		g.setScanWindow(samples.Sim.DtSamples)
		return true
	}

	return false
}

func (g *DTGraph) setScanWindow(lanes *samples.Samples) {
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

func (g *DTGraph) dataAccessor() (x, y float64, c color.Color, more int) {
	// Return values until we reach the end of the current active lane.

	if g.scanIdx >= g.scanEnd {
		// No more samples in the lane.
		g.setScanWindow(samples.Sim.DtSamples)
		return 0, 0, nil, 0
	}

	sample := g.activeLane.Values[g.scanIdx]

	g.scanIdx++

	if sample.Value != nil {
		return sample.Time, sample.Value.(float64), g.lineColor, 1
	}

	return 0, 0, nil, 0

}
