package gui

import (
	"image/color"

	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

type Button struct {
	baseButton

	selected bool

	selectedColor   color.RGBA
	unselectedColor color.RGBA

	// Flash properties
	flashable bool
	delay     int
	delayCnt  int
}

func NewButton(parent IWidget, width, height int) IWidget {
	gb := new(Button)
	gb.initialize(parent, width, height)

	gb.selectedColor = color.RGBA{200, 127, 127, 127}
	gb.unselectedColor = color.RGBA{127, 127, 127, 127}

	gb.ev.Source = "Button"
	gb.ev.Target = "App"
	gb.ev.Action = "Command"

	gb.delay = 2
	gb.delayCnt = gb.delay
	return gb
}

func (gb *Button) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	if eventType != events.MouseButton { // button event
		return false, -1
	}

	if !gb.isInside(x, y) {
		return false, -1
	}

	gb.delayCnt = 0

	comm.MsgBus.SendEvent(gb.ev)
	comm.MsgBus.Send2("Gui", "Button", "Changed", "", "", "")

	return true, gb.id
}

func (gb *Button) Select() {
	gb.selected = true
}

func (gb *Button) UnSelect() {
	gb.selected = false
}

func (gb *Button) DrawAt(x, y int32) {
	gb.SetPos(x, y)
	gb.Draw()
}

func (gb *Button) Draw() {
	if gb.flashable {
		if gb.delayCnt < gb.delay {
			gb.DC.SetColor(gb.selectedColor)
		} else {
			gb.DC.SetColor(gb.unselectedColor)
		}
		gb.delayCnt++
	} else {
		if gb.selected {
			gb.DC.SetColor(gb.selectedColor)
		} else {
			gb.DC.SetColor(gb.unselectedColor)
		}
	}

	gb.DC.DrawRectangle(float64(gb.x), float64(gb.y), float64(gb.Rect.W), float64(gb.Rect.H))
	gb.DC.Fill()

	gb.DC.SetColor(gb.borderColor)
	gb.DC.DrawRectangle(float64(gb.x), float64(gb.y), float64(gb.Rect.W), float64(gb.Rect.H))
	gb.DC.Stroke()

	gb.DC.SetColor(gb.labelColor)
	gb.DC.DrawString(gb.label, float64(gb.x+gb.Rect.W/2)-gb.textWidth/2, float64(gb.y+gb.Rect.H/2)+gb.textHeight/2)
}
