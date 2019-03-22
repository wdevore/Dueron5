package gui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

// Misc Panel

type Km1PanelWidget struct {
	basePanel

	fields *ValueGroup
}

func NewKm1PanelWidget(startID int, renderer *sdl.Renderer, texture *sdl.Texture, width, height int) (widget IWidget, id int) {
	pw := new(Km1PanelWidget)
	pw.initialize(pw, renderer, texture, width, height)
	pw.canHide = true
	pw.id = startID
	iDs := pw.build()

	return pw, iDs
}

func (pw *Km1PanelWidget) Listen(msg *comm.MessageEvent) {
	switch msg.Target {
	case "Panel":
		switch msg.Action {
		case "Set":
			_, widget := pw.fields.values.Find(func(id int, v interface{}) bool {
				return v.(*ValueButton).ev.Field == msg.Field
			})
			widget.(*ValueButton).SetValue(msg.Value)
			break
		}
		break
	}
}

// override's baseWidget and implements IWidget
func (pw *Km1PanelWidget) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	if !pw.visible {
		return false, -1
	}

	handled, id = pw.fields.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	return false, -1
}

func (pw *Km1PanelWidget) Refresh() {
	it := pw.fields.Iterator()
	for it.Next() {
		btn := it.Value().(*ValueButton)
		btn.Refresh()
	}
}

func (pw *Km1PanelWidget) DrawAt(x, y int32) {
	pw.SetPos(x, y)
	pw.Draw()
}

func (pw *Km1PanelWidget) Draw() {
	if !pw.visible {
		return
	}

	pw.DC.SetColor(pw.backgroundColor)
	pw.DC.Clear()

	// Draw child widgets
	it := pw.fields.Iterator()
	for it.Next() {
		btn := it.Value().(*ValueButton)
		btn.Draw()
	}

	pw.update()
}

func (pw *Km1PanelWidget) build() int {
	pw.fields = NewValueGroup()

	iDs := pw.id
	btnXPos := int32(0)
	btnYPos := int32(0)

	wigBut := NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*ValueButton)
	btn.ev.Message = "Simulation"
	btn.ev.Field = "Duration"
	btn.SetLabel("Duration")
	v := deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Simulation"
	btn.ev.Field = "TimeStep"
	btn.Min = 0
	btn.Max = 3600.0 * 10         // = 10 seconds
	btn.SetLabel("Time Step(us)") // µs
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Simulation"
	btn.ev.Field = "StimulusScaler"
	btn.Min = 0
	btn.Max = 100
	btn.SetLabel("Stim Scale")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Simulation"
	btn.ev.Field = "Hertz"
	btn.Min = 0
	btn.Max = 1000
	btn.SetLabel("Hertz")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Simulation"
	btn.ev.Field = "Firing_Rate"
	btn.Min = 0
	btn.Max = 50.0
	btn.SetLabel("Firing Rate")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	iDs++

	return iDs
}
