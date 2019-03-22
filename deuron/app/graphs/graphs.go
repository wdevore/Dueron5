package graphs

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/gui"
	"github.com/wdevore/Deuron5/simulation/samples"
)

const (
	SelectBarWidth = 4
)

// SeriesAccessor provides access to series data
type SeriesAccessor func() (x, y float64, c color.Color, state int)

var ran = rand.New(rand.NewSource(1963))

type IGraph interface {
	Destroy()
	Prep()
	Draw()
	SetSeries(SeriesAccessor)
	Check() bool
	MarkDirty(dirty bool)
	SetName(string)
	Name() string
}

type BaseGraph struct {
	name       string
	titleWidth float64
	Bounds     sdl.Rect

	pixels *image.RGBA

	borderOffsetX float64
	borderOffsetY float64

	upperX float64
	upperY float64

	backgroundColor color.RGBA
	borderColor     color.RGBA

	mouseVertColor   color.RGBA
	mouseTextColor   color.RGBA
	titleTextColor   color.RGBA
	maxTextColor     color.RGBA
	zeroLineColor    color.RGBA
	selectedBarColor color.RGBA

	// Mouse in view-space (aka world-space)
	WVx, WVy int32
	// Mouse coords in graph-space (aka local-space)
	MVx, MVy int32
	// Mouse coords in unit-space
	Ux, Uy float64
	// coords in graph-space
	GVx, GVy float64

	renderer *sdl.Renderer
	texture  *sdl.Texture

	scanIdx   int
	scanStart int
	scanEnd   int

	dirty bool

	selected bool

	activeLane *samples.SamplesLane
}

func (bg *BaseGraph) Initialize(rect sdl.Rect, dc *gg.Context) {
	bg.backgroundColor = color.RGBA{48, 48, 48, 255}
	bg.borderColor = color.RGBA{127, 127, 127, 255}

	bg.mouseVertColor = color.RGBA{200, 200, 200, 255}
	bg.mouseTextColor = color.RGBA{255, 127, 50, 255}
	bg.maxTextColor = color.RGBA{127, 255, 127, 255}
	bg.titleTextColor = color.RGBA{255, 255, 255, 255}
	bg.zeroLineColor = color.RGBA{127, 127, 255, 127}
	// bg.zeroLineColor = color.RGBA{255, 255, 255, 255}
	bg.selectedBarColor = color.RGBA{127, 127, 255, 255}

	bg.borderOffsetX = 4.0
	bg.borderOffsetY = 4.0

	bg.pixels = dc.Image().(*image.RGBA)

	bg.upperX = float64(rect.W) - bg.borderOffsetX*2.0
	bg.upperY = float64(rect.H) - bg.borderOffsetY*2.0

}

func (bg *BaseGraph) SetName(name string) {
	bg.name = name
}

func (bg *BaseGraph) Name() string {
	return bg.name
}

func (bg *BaseGraph) SetGraphics(renderer *sdl.Renderer, texture *sdl.Texture) {
	bg.renderer = renderer
	bg.texture = texture
}

func (bg *BaseGraph) MarkDirty(dirty bool) {
	bg.dirty = dirty
}

func (bg *BaseGraph) Prep() {
}

func (bg *BaseGraph) Draw(rect sdl.Rect, dc *gg.Context) {
	dc.Identity()

	dc.SetColor(bg.backgroundColor)
	dc.Clear()

	if bg.selected {
		bg.drawSelectBar(rect.H, dc)
	}

	// Render the graph onto the image pixels

	dc.InvertY()

	dc.Translate(bg.borderOffsetX, bg.borderOffsetY)

	dc.SetLineWidth(1.0)

	// Box borders
	bg.drawBoxBorder(dc)
}

func (bg *BaseGraph) postDraw(rect sdl.Rect) {
	// -------------------------------------------
	// Blit pixels
	// -------------------------------------------
	bg.texture.Update(&rect, bg.pixels.Pix, bg.pixels.Stride)

	// Now copy the texture onto the target (aka the display)
	bg.renderer.Copy(bg.texture, &rect, &rect)
}

func (bg *BaseGraph) drawSelectBar(height int32, dc *gg.Context) {
	dc.SetColor(bg.selectedBarColor)
	dc.SetLineWidth(SelectBarWidth)
	dc.MoveTo(SelectBarWidth, 0.0)
	dc.LineTo(SelectBarWidth, float64(height))
	dc.Stroke()
}

func (bg *BaseGraph) drawZeroLine(rect sdl.Rect, dc *gg.Context) {
	dc.Identity()
	dc.SetLineWidth(1.0)
	dc.Translate(bg.borderOffsetX, bg.borderOffsetY)

	// we need to map "backwards" in order to position the vertical bar
	// withing the scan window. As usual we map to unit-space so we can map
	// to whatever destination space needed, in this case we map back to
	// window-space. We also need to truncate from float to int so the bar
	// jumps from "t" to "t".
	uspY := bg.Linear(bg.activeLane.Max, bg.activeLane.Min, float64(0))

	// Map from unit-space to window-space
	winY := float64(int(bg.Lerp(0, bg.upperY, uspY)))
	// fmt.Printf("%f\n", winY)
	dc.SetColor(bg.zeroLineColor)
	dc.MoveTo(0, winY)
	dc.LineTo(float64(rect.W)-bg.borderOffsetX*2, winY)
	dc.Stroke()
}

func (bg *BaseGraph) drawVerticalTimeBar(rect sdl.Rect, dc *gg.Context) float64 {
	dc.Identity()
	dc.Translate(bg.borderOffsetX, bg.borderOffsetY)
	dc.SetLineWidth(1.0)

	// we need to map "backwards" in order to position the vertical bar
	// withing the scan window. As usual we map to unit-space so we can map
	// to whatever destination space needed, in this case we map back to
	// window-space. We also need to truncate from float to int so the bar
	// jumps from "t" to "t".
	uspX := bg.Linear(float64(bg.scanStart), float64(bg.scanEnd), float64(int(bg.GVx)))

	// Map from unit-space to window-space
	winX := bg.Lerp(0, bg.upperX, uspX)

	dc.SetColor(bg.mouseVertColor)
	dc.MoveTo(float64(int(winX)), 0)
	dc.LineTo(float64(int(winX)), float64(rect.H)-bg.borderOffsetY*2)
	dc.Stroke()

	return winX
}

func (bg *BaseGraph) drawTitle(rect sdl.Rect, dc *gg.Context) {
	dc.Identity()
	dc.Translate(bg.borderOffsetX, bg.borderOffsetY)
	dc.SetColor(bg.titleTextColor)
	bg.titleWidth, _ = dc.MeasureString(bg.Name())
	upperYInv := float64(rect.H) - bg.upperY + bg.borderOffsetY
	upperXInv := float64(rect.W) - bg.titleWidth - bg.borderOffsetX*2.0
	dc.DrawString(bg.name, upperXInv, upperYInv)
}

func (bg *BaseGraph) drawMinMax(min, max float64, rect sdl.Rect, dc *gg.Context) {
	dc.Identity()
	dc.Translate(bg.borderOffsetX, bg.borderOffsetY)
	dc.SetColor(bg.maxTextColor)
	bg.titleWidth, _ = dc.MeasureString(bg.Name())
	upperYInv := float64(rect.H) - bg.upperY + bg.borderOffsetY + 15
	dc.DrawString(fmt.Sprintf("[%0.3f, %0.3f]", min, max), 5, upperYInv)
}

func (bg *BaseGraph) drawMouseInfo(winX float64, rect sdl.Rect, dc *gg.Context) {
	dc.Identity()
	dc.Translate(bg.borderOffsetX, bg.borderOffsetY)
	dc.SetColor(bg.mouseTextColor)

	upperYInv := float64(rect.H) - bg.upperY + bg.borderOffsetY
	dc.DrawString(fmt.Sprintf("(%d, %0.3f)", int(bg.GVx), bg.GVy), winX+1, upperYInv)

	// g.DC.DrawString(fmt.Sprintf("L(%0.2f, %0.2f)", g.GVx, g.GVy), 5, upperYInv)
	// g.DC.DrawString(fmt.Sprintf("%d", int(g.GVx)), winX+1, upperYInv)
	// g.DC.DrawString(fmt.Sprintf("L(%d, %d)", int(g.GVx), int(g.GVy)), 5, upperYInv)
	// dc.DrawString(fmt.Sprintf("M(%d, %d)", bg.MVx, bg.MVy), 5, upperYInv+15)

	// g.DC.DrawString(fmt.Sprintf("W(%d, %d)", g.WVx, g.WVy), 5, upperYInv)
	// g.DC.DrawString(fmt.Sprintf("U(%0.2f, %0.2f)", g.Ux, g.Uy), 5, upperYInv+30)
}

func (bg *BaseGraph) drawBoxBorder(dc *gg.Context) {
	dc.SetColor(bg.borderColor)

	// Left border
	dc.MoveTo(0.0, 0.0)
	dc.LineTo(0.0, bg.upperY)

	// Top border
	dc.LineTo(bg.upperX, bg.upperY)

	// Right border
	dc.LineTo(bg.upperX, 0.0)

	dc.ClosePath()

	dc.Stroke()
}

func (bg *BaseGraph) preHandle(msg *comm.MessageEvent) bool {
	if !bg.selected {
		return true
	}

	return false
}

func (bg *BaseGraph) handleMotion(vx, vy int32, widget gui.BaseWidget, inside bool) bool {
	if inside && bg.selected {
		bg.WVx = vx
		bg.WVy = vy

		px, py := widget.Position()
		bg.MapToGraphSpace(widget.Rect, px, py)

		comm.MsgBus.Send2("Gui", "Data", "Changed", "", "", "")

		return true
	}

	return false
}

func (bg *BaseGraph) handleScroll(msg *comm.MessageEvent, sams *samples.Samples, px, py int32, rect sdl.Rect) bool {
	switch msg.Target {
	case "Key":
		switch msg.Action {
		case "Left":
			switch msg.Message {
			case "ScrollLeft":
				// Scroll window left until reach min 0. The window width
				// doesn't change.
				bg.scrollLeft("PspGraph", sams)
				bg.MapToGraphSpace(rect, px, py)
				return true
			}
			break
		case "Right":
			switch msg.Message {
			case "ScrollRight":
				bg.scrollRight("PspGraph", sams)
				bg.MapToGraphSpace(rect, px, py)
				return true
			}
			break
		}
		break
	}

	return false
}

func (bg *BaseGraph) handleRange(msg *comm.MessageEvent, sams *samples.Samples) bool {
	switch msg.Target {
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
				sams.SetRangeStart(int(start))
				return true
			case "Lane_End":
				// Update sample property
				end, err := strconv.ParseFloat(msg.Value, 64)
				if err != nil {
					fmt.Printf("Error Start: %v\n", err)
				}
				sams.SetRangeEnd(int(end))
				return true
			}
			break
		}
		break
	}

	return false
}

func (bg *BaseGraph) scrollLeft(source string, sams *samples.Samples) {
	start, end := sams.GetRange()

	if start > 0 {
		inc := deuron.SimModel.GetFloat("Inc/Dec")

		start = start - int(inc)
		if start >= 0 {
			end = end - int(inc)
			comm.MsgBus.Send3(source, "Model", "Set", "", "", "Lane_Start", fmt.Sprintf("%d", start))
			comm.MsgBus.Send3(source, "Model", "Set", "", "", "Lane_End", fmt.Sprintf("%d", end))
		}
	} else {
		// fmt.Println("Can't scroll left")
	}
}

func (bg *BaseGraph) scrollRight(source string, sams *samples.Samples) {
	start, end := sams.GetRange()

	if end <= sams.Size() { //  or <
		inc := deuron.SimModel.GetFloat("Inc/Dec")
		end = end + int(inc)
		if end <= sams.Size() {
			start = start + int(inc)
			comm.MsgBus.Send3(source, "Model", "Set", "", "", "Lane_Start", fmt.Sprintf("%d", start))
			comm.MsgBus.Send3(source, "Model", "Set", "", "", "Lane_End", fmt.Sprintf("%d", end))
		}
	} else {
		// fmt.Println("Can't scroll right")
	}
}

// -----------------------------------------------------------------------
// Space Mappings
// -----------------------------------------------------------------------
func (bg *BaseGraph) MapToGraphSpace(rect sdl.Rect, px, py int32) {
	// Map world-space to local-space
	bg.MVx = bg.WVx - px
	bg.MVy = bg.WVy - py
	// Map local mouse-space to unit-space
	bg.Ux = bg.Linear(bg.borderOffsetX, float64(rect.W)-bg.borderOffsetX, float64(bg.MVx))
	bg.Uy = bg.Linear(bg.borderOffsetY, float64(rect.H)-bg.borderOffsetY, float64(bg.MVy))

	// Map from unit-space to window-space
	// 0 is the min because we have already translated the origin.
	// g.LVx = g.Lerp(0, float64(g.Rect.W), g.Ux)
	// g.LVy = g.Lerp(0, float64(g.Rect.H), g.Uy)

	bg.GVx = bg.Lerp(float64(bg.scanStart), float64(bg.scanEnd), bg.Ux)
	if bg.activeLane != nil {
		// Y is inverted so we map from max to min
		bg.GVy = bg.Lerp(bg.activeLane.Max, bg.activeLane.Min, bg.Uy)
	} else {
		bg.GVy = bg.Lerp(0, float64(rect.H), bg.Uy)
	}
}

// Linear returns 0->1 for a "value" between min and max.
// Generally used to map from view-space to unit-space
func (bg *BaseGraph) Linear(min, max, value float64) float64 {
	if min < 0 {
		return (value - max) / (min - max)
	}

	return (value - min) / (max - min)
}

// Lerp returns a the value between min and max given t = 0->1
// Call Linear first to get "t".
// Generally used to map from unit-space to window-space
func (bg *BaseGraph) Lerp(min, max, t float64) float64 {
	return min*(1.0-t) + max*t
}
