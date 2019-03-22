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

// Renders the Surge value for a particular synapse.

type NeuronAPGraph struct {
	BaseGraph
	gui.BaseWidget

	lineColor    color.RGBA
	altLineColor color.RGBA
	even         bool

	prevWinX, prevWinY float64

	graphIt sll.Iterator
}

func NewNeuronAPGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) IGraph {
	g := new(NeuronAPGraph)
	g.BaseWidget.Initialize(g, width, height)
	g.BaseGraph.Initialize(g.Rect, g.DC)
	g.SetGraphics(renderer, texture)

	g.lineColor = color.RGBA{255, 127, 0, 255}
	g.altLineColor = color.RGBA{255, 160, 0, 255}

	return g
}

func (g *NeuronAPGraph) Listen(msg *comm.MessageEvent) {
	preHandled := g.preHandle(msg)

	if preHandled {
		return
	}

	samples := samples.Sim.NeuronAPSamples

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

func (g *NeuronAPGraph) Handle(vx, vy int32, eventType events.MouseEventType) (handled bool, id int) {
	inside := gui.PointInside(vx, vy, g.Rect.X, g.Rect.Y, g.Rect.W, g.Rect.H)

	switch eventType {
	case events.MouseButton:
		if inside {
			g.selected = !g.selected
			if g.selected {
				// Send message to app.go
				comm.MsgBus.Send2("NeuronAPGraph", "Graph", "Selected", "Surge", fmt.Sprintf("%d", g.ID()), "")
				// Update gui Range fields the sample's range.
				samples := samples.Sim.NeuronAPSamples
				start, end := samples.GetRange()

				// Send message to panel including the panel's id.
				comm.MsgBus.Send3("NeuronAPGraph", "Model", "Set", "", "", "Lane_Start", fmt.Sprintf("%d", start))
				comm.MsgBus.Send3("NeuronAPGraph", "Model", "Set", "", "", "Lane_End", fmt.Sprintf("%d", end))
			} else {
				comm.MsgBus.Send2("NeuronAPGraph", "Graph", "UnSelected", "Surge", fmt.Sprintf("%d", g.ID()), "")
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

func (g *NeuronAPGraph) SetSeries(accessor SeriesAccessor) {
}

// Destroy release resources
func (g *NeuronAPGraph) Destroy() {
}

func (g *NeuronAPGraph) Prep() {
}

// DrawAt renders graph to texture
func (g *NeuronAPGraph) Draw() {
	g.BaseGraph.Draw(g.Rect, g.DC)

	if g.activeLane.Min < 0 {
		g.drawZeroLine(g.Rect, g.DC)
	}

	// Draw ticks
	g.drawXticks()

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

		g.prevWinX = winX
		g.prevWinY = winY

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

func (g *NeuronAPGraph) drawXticks() {
	// Clip to range window
	// duration := deuron.SimModel.GetFloat("Samples")

	// for x := float64(0); x < duration; x++ {
	// 	xi := int(x) % 100
	// 	if xi == 0 {
	// 		// Draw tick
	// 		uspX := g.Linear(0, duration, float64(xi))
	// 		winX := g.Lerp(0, g.upperX, uspX)
	// 	}
	// }
}

func (g *NeuronAPGraph) Check() bool {
	if samples.Sim.NeuronAPSamples == nil {
		return false
	}

	lanes := samples.Sim.NeuronAPSamples.GetLanes()

	g.graphIt = lanes.Iterator()

	if g.graphIt.First() {
		lane, _ := lanes.Get(0)
		g.activeLane = lane.(*samples.SamplesLane)
		g.setScanWindow(samples.Sim.NeuronAPSamples)
		return true
	}

	return false
}

func (g *NeuronAPGraph) setScanWindow(lanes *samples.Samples) {
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

func (g *NeuronAPGraph) dataAccessor() (x, y float64, c color.Color, more int) {
	// Return values until we reach the end of the current active lane.

	if g.scanIdx >= g.scanEnd {
		// No more samples in the lane.
		g.setScanWindow(samples.Sim.NeuronAPSamples)
		return 0, 0, nil, 0
	}

	sample := g.activeLane.Values[g.scanIdx]

	g.scanIdx++

	if sample.Value != nil {
		return sample.Time, sample.Value.(float64), g.lineColor, 1
	}

	return 0, 0, nil, 0

}
