package gui

import (
	"image/color"

	"github.com/wdevore/Deuron5/deuron/app/events"

	"github.com/wdevore/Deuron5/deuron/app/comm"
)

type ToggleButton struct {
	baseButton

	selected bool

	selectedColor   color.RGBA
	unselectedColor color.RGBA
}

func NewToggleButton(parent IWidget, width, height int) IWidget {
	gb := new(ToggleButton)
	gb.initialize(parent, width, height)

	gb.selectedColor = color.RGBA{127, 200, 127, 127}
	gb.unselectedColor = color.RGBA{100, 100, 64, 127}

	gb.ev.Source = "ToggleButton"
	gb.ev.Target = "Model"
	gb.ev.Action = "Toggle"

	return gb
}

func (gb *ToggleButton) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	if eventType != events.MouseButton { // button event
		return false, -1
	}

	if !gb.isInside(x, y) {
		return false, -1
	}

	gb.selected = !gb.selected

	comm.MsgBus.SendEvent(gb.ev)
	comm.MsgBus.Send2("Gui", "ToggleButton", "Changed", "", "", "")

	return true, gb.id
}

func (gb *ToggleButton) Select() {
	gb.selected = true
}

func (gb *ToggleButton) UnSelect() {
	gb.selected = false
}

func (gb *ToggleButton) DrawAt(x, y int32) {
	gb.SetPos(x, y)
	gb.Draw()
}

func (gb *ToggleButton) Draw() {
	if gb.selected {
		gb.DC.SetColor(gb.selectedColor)
	} else {
		gb.DC.SetColor(gb.unselectedColor)
	}
	gb.DC.DrawRectangle(float64(gb.x), float64(gb.y), float64(gb.Rect.W), float64(gb.Rect.H))
	gb.DC.Fill()

	gb.DC.SetColor(gb.borderColor)
	gb.DC.DrawRectangle(float64(gb.x), float64(gb.y), float64(gb.Rect.W), float64(gb.Rect.H))
	gb.DC.Stroke()

	gb.DC.SetColor(gb.labelColor)
	gb.DC.DrawString(gb.label, float64(gb.x+gb.Rect.W/2)-gb.textWidth/2, float64(gb.y+gb.Rect.H/2)+gb.textHeight/2)
}
