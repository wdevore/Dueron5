package gui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

type Km2PanelWidget struct {
	basePanel

	stimGroup *ButtonGroup
}

func NewKm2PanelWidget(startID int, renderer *sdl.Renderer, texture *sdl.Texture, width, height int) (widget IWidget, id int) {
	pw := new(Km2PanelWidget)
	pw.initialize(pw, renderer, texture, width, height)
	pw.canHide = true
	pw.id = startID
	iDs := pw.build()

	return pw, iDs
}

func (pw *Km2PanelWidget) Listen(msg *comm.MessageEvent) {
}

// override's baseWidget and implements IWidget
func (pw *Km2PanelWidget) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	if !pw.visible {
		return false, -1
	}

	handled, id = pw.stimGroup.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	return false, -1
}

func (pw *Km2PanelWidget) DrawAt(x, y int32) {
	pw.SetPos(x, y)
	pw.Draw()
}

func (pw *Km2PanelWidget) Draw() {
	if !pw.visible {
		return
	}

	pw.DC.SetColor(pw.backgroundColor)
	pw.DC.Clear()

	// Draw child widgets
	it := pw.stimGroup.Iterator()
	for it.Next() {
		btn := it.Value().(*Button)
		btn.Draw()
	}

	pw.update()
}

func (pw *Km2PanelWidget) build() int {

	iDs := pw.id
	btnXPos := int32(0)
	btnYPos := int32(0)

	pw.stimGroup = NewButtonGroup()

	wigBut := NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Stimulus"
	btn.ev.Value = "stim_1"
	btn.ev.Message = "Simulation"
	btn.SetLabel("Stim 01")
	iDs++
	wigBut.SetID(iDs)
	pw.stimGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Stimulus"
	btn.ev.Value = "stim_2"
	btn.ev.Message = "Simulation"
	btn.SetLabel("Stim 02")
	iDs++
	wigBut.SetID(iDs)
	pw.stimGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Stimulus"
	btn.ev.Value = "stim_3"
	btn.ev.Message = "Simulation"
	btn.SetLabel("Stim 03")
	iDs++
	wigBut.SetID(iDs)
	pw.stimGroup.AddButton(wigBut)

	iDs++

	return iDs
}
