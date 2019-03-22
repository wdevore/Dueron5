package gui

import (
	"fmt"
	"image/color"

	"github.com/wdevore/Deuron5/deuron/app/events"

	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
)

const (
	SubButtonWidth  = 25
	SubButtonHeight = 25
)

// Value buttons have three sub-buttons for inc/dec and edit.

type ValueButton struct {
	baseButton

	valueWidth  float64
	valueHeight float64

	valueColor    color.RGBA
	incDecColor   color.RGBA
	editTextColor color.RGBA

	// The "Escape" key cancels edit.
	editMode bool

	originalValue string

	// Constraints
	Constrained bool
	Min, Max    float64
}

func NewValueButton(parent IWidget, width, height int) IWidget {
	gb := new(ValueButton)
	gb.initialize(parent, width, height)

	gb.valueColor = color.RGBA{255, 255, 8, 255}
	gb.incDecColor = color.RGBA{255, 32, 32, 127}
	gb.editTextColor = color.RGBA{64, 255, 64, 255}

	gb.ev.Source = "ValueButton"

	comm.MsgBus.Subscribe(gb)

	gb.editMode = false
	return gb
}

func (gb *ValueButton) Listen(msg *comm.MessageEvent) {
	if !(msg.Field == gb.ev.Field || msg.ID == fmt.Sprintf("%d", gb.id)) {
		return
	}

	switch msg.Target {
	// These messages are sent from keymaps.go
	case "Control":
		switch msg.Action {
		case "Edit":
			switch msg.Message {
			case "Cancel":
				gb.editMode = false
				gb.SetValue(gb.originalValue)
				return
			case "Done":
				gb.editMode = false
				comm.MsgBus.Send3(gb.ev.Source, "Model", "Set", gb.ev.Message, gb.ev.ID, gb.ev.Field, gb.ev.Value)
			case "Current":
				gb.SetValue(msg.Value)
				comm.MsgBus.Send2("Gui", "Data", "Changed", "", "", "")
				return
			}
			return
		}
		break
	case "Data":
		switch msg.Action {
		case "Changed":
			// Update Gui
			gb.SetValue(msg.Value)
			// gb.simUpdate()
			return
		}
		break
	}
}

func (gb *ValueButton) SetValue(value string) {
	gb.ev.Value = value
	gb.valueWidth, gb.valueHeight = gb.DC.MeasureString(value)
}

func (gb *ValueButton) Refresh() {
	value := deuron.SimModel.GetFloatAsString(gb.ev.Field)
	// fmt.Printf("%s\n", value)
	gb.SetValue(value)
}

func (gb *ValueButton) Handle(vx, vy int32, eventType events.MouseEventType) (handled bool, id int) {
	if eventType != events.MouseButton { // button event
		return false, -1
	}

	if !gb.isInside(vx, vy) {
		return false, -1
	}

	// Check for sub-button hits.
	px, py := gb.parent.Position()
	lx, ly := ViewSpaceToLocal(vx, vy, px, py)

	// Inc
	if PointInside(lx, ly, gb.x+gb.Rect.W-SubButtonWidth, gb.y, SubButtonWidth, SubButtonHeight) {
		// Notify model of a new value.
		inc := deuron.SimModel.GetFloat("Inc/Dec")
		fValue := deuron.SimModel.GetFloat(gb.ev.Field)
		fValue = fValue + inc

		if !gb.Constrained && fValue <= gb.Max {
			comm.MsgBus.Send3(gb.ev.Source, "Model", "Set", gb.ev.Message, gb.ev.ID, gb.ev.Field, fmt.Sprintf("%0.3f", fValue))
		}
		return true, gb.id
	}

	// Dec
	if PointInside(lx, ly, gb.x+gb.Rect.W-SubButtonWidth, gb.y+SubButtonHeight, SubButtonWidth, SubButtonHeight) {
		inc := deuron.SimModel.GetFloat("Inc/Dec")
		fValue := deuron.SimModel.GetFloat(gb.ev.Field)
		fValue = fValue - inc

		if !gb.Constrained && fValue >= gb.Min {
			comm.MsgBus.Send3(gb.ev.Source, "Model", "Set", gb.ev.Message, gb.ev.ID, gb.ev.Field, fmt.Sprintf("%0.3f", fValue))
		}
		return true, gb.id
	}

	// Edit
	if PointInside(lx, ly, gb.x+gb.Rect.W-SubButtonWidth*2, gb.y+SubButtonHeight, SubButtonWidth, SubButtonHeight) {
		// Put widget into the "Edit" state. The "Escape" key cancels edit.
		gb.editMode = true
		gb.originalValue = deuron.SimModel.GetFloatAsString(gb.ev.Field)

		comm.MsgBus.Send3(gb.ev.Source, "App", "Edit", "Start", gb.ev.ID, gb.ev.Field, gb.originalValue)

		return true, gb.id
	}

	return true, gb.id
}

func (gb *ValueButton) DrawAt(x, y int32) {
	gb.SetPos(x, y)
	gb.Draw()
}

func (gb *ValueButton) Draw() {
	// Draw Inc sub-button in top-right corner
	gb.DC.SetColor(gb.incDecColor)
	gb.DC.DrawRectangle(float64(gb.x+gb.Rect.W-SubButtonWidth), float64(gb.y), float64(SubButtonWidth), float64(SubButtonHeight*2))
	// Edit button background
	gb.DC.DrawRectangle(float64(gb.x+gb.Rect.W-SubButtonWidth*2), float64(gb.y+SubButtonHeight), float64(SubButtonWidth), float64(SubButtonHeight))
	gb.DC.Fill()

	gb.DC.SetColor(gb.borderColor)
	// Inc sub-button outline
	gb.DC.DrawRectangle(float64(gb.x+gb.Rect.W-SubButtonWidth), float64(gb.y), float64(SubButtonWidth), float64(gb.Rect.H/2))
	// Dec sub-button outline
	gb.DC.DrawRectangle(float64(gb.x+gb.Rect.W-SubButtonWidth), float64(gb.y+SubButtonHeight), float64(SubButtonWidth), float64(SubButtonHeight))
	// Edit sub-button outline
	gb.DC.DrawRectangle(float64(gb.x+gb.Rect.W-SubButtonWidth*2), float64(gb.y+SubButtonHeight), float64(SubButtonWidth), float64(SubButtonHeight))
	// Value button outline
	gb.DC.DrawRectangle(float64(gb.x), float64(gb.y), float64(gb.Rect.W), float64(gb.Rect.H))
	gb.DC.Stroke()

	gb.DC.SetColor(gb.labelColor)
	gb.DC.DrawString("E", float64(gb.x+gb.Rect.W-SubButtonWidth*2)+(SubButtonWidth/2), float64(gb.y+gb.Rect.H/2)+(SubButtonHeight/2))
	// gb.DC.DrawString(gb.label, float64(gb.x+gb.Rect.W/2)-gb.textWidth/2, float64(gb.y)+gb.textHeight)
	gb.DC.DrawString(gb.label, float64(gb.x+5), float64(gb.y)+gb.textHeight)

	if gb.editMode {
		gb.DC.SetColor(gb.editTextColor)
	} else {
		gb.DC.SetColor(gb.valueColor)
	}
	gb.DC.DrawString(gb.ev.Value, float64(gb.x+10), float64(gb.y+gb.Rect.H/2)+gb.textHeight/2)
}

// func (gb *ValueButton) simUpdate() {
// 	autoRun := deuron.SimModel.GetFloat("AutoRunPause")
// 	if autoRun == 1 {
// 		comm.MsgBus.Send(gb.ev.Source, "App", "Command", "RunPause", "")
// 	}
// }
