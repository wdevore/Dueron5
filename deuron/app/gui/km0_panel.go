package gui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

type Km0PanelWidget struct {
	basePanel

	fields *ValueGroup
}

func NewKm0PanelWidget(startID int, renderer *sdl.Renderer, texture *sdl.Texture, width, height int) (widget IWidget, id int) {
	pw := new(Km0PanelWidget)
	pw.initialize(pw, renderer, texture, width, height)
	pw.canHide = true
	pw.id = startID
	iDs := pw.build()

	return pw, iDs
}

func (pw *Km0PanelWidget) Listen(msg *comm.MessageEvent) {
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
func (pw *Km0PanelWidget) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	if !pw.visible {
		return false, -1
	}

	handled, id = pw.fields.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	return false, -1
}

func (pw *Km0PanelWidget) Refresh() {
	it := pw.fields.Iterator()
	for it.Next() {
		btn := it.Value().(*ValueButton)
		btn.Refresh()
	}
}

func (pw *Km0PanelWidget) DrawAt(x, y int32) {
	pw.SetPos(x, y)
	pw.Draw()
}

func (pw *Km0PanelWidget) Draw() {
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

func (pw *Km0PanelWidget) build() int {
	pw.fields = NewValueGroup()

	iDs := pw.id
	btnXPos := int32(0)
	btnYPos := int32(0)

	iDs = pw.build_row1(iDs, btnXPos, btnYPos)

	// Next Row
	btnXPos = 0
	btnYPos = btnYPos + DefaultValueButtonHeight
	iDs = pw.build_row2(iDs, btnXPos, btnYPos)

	btnXPos = 0
	btnYPos = btnYPos + DefaultValueButtonHeight
	iDs = pw.build_row3(iDs, btnXPos, btnYPos)

	btnXPos = 0
	btnYPos = btnYPos + DefaultValueButtonHeight
	iDs = pw.build_row4(iDs, btnXPos, btnYPos)

	btnXPos = 0
	btnYPos = btnYPos + DefaultValueButtonHeight
	iDs = pw.build_row5(iDs, btnXPos, btnYPos)

	return iDs
}

func (pw *Km0PanelWidget) build_row1(iDs int, btnXPos, btnYPos int32) int {
	wigBut := NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*ValueButton)
	btn.ev.Field = "Range_Start"
	btn.SetLabel("Range start")
	v := deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.Max = deuron.SimModel.GetFloat("Range_End")
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Field = "Range_End"
	btn.SetLabel("Range end")
	btn.Max = deuron.SimModel.GetFloat("Samples")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Field = "Lane_Start"
	btn.SetLabel("Lane Start")
	btn.Max = 10000000.0
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Field = "Lane_End"
	btn.SetLabel("Lane End")
	btn.Max = deuron.SimModel.GetFloat("Samples")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	iDs++

	return iDs
}

// func (pw *Km0PanelWidget) build_row2(iDs int, btnXPos, btnYPos int32) int {

// 	// btnXPos = btnXPos + DefaultButtonWidth + 50
// 	// wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
// 	// wigBut.SetPos(btnXPos, btnYPos)
// 	// btn = wigBut.(*ValueButton)
// 	// btn.ev.Field = "Poisson_min"
// 	// btn.Max = 5000.0
// 	// btn.SetLabel("Poisson min")
// 	// v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
// 	// btn.SetValue(v)
// 	// iDs++
// 	// wigBut.SetID(iDs)
// 	// pw.fields.AddWidget(wigBut)

// 	iDs++

// 	return iDs
// }

func (pw *Km0PanelWidget) build_row2(iDs int, btnXPos, btnYPos int32) int {
	wigBut := NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*ValueButton)
	btn.ev.Field = "Active_Synapse"
	btn.Max = deuron.SimModel.GetFloat("Synapse_Count") - 1
	btn.SetLabel("Synapse")
	v := deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Neuron"
	btn.ev.Field = "threshold"
	btn.Max = 1000.0
	btn.SetLabel("Threshold")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	iDs++

	return iDs
}

func (pw *Km0PanelWidget) build_row3(iDs int, btnXPos, btnYPos int32) int {
	wigBut := NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "taoP"
	btn.Max = 1000.0
	btn.SetLabel("taoP")
	v := deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "taoN"
	btn.Max = 1000.0
	btn.SetLabel("taoN")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "taoI"
	btn.Max = 1000.0
	btn.SetLabel("taoI")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "amb"
	btn.Max = 1000.0
	btn.SetLabel("amb")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "ama"
	btn.Max = 1000.0
	btn.SetLabel("ama")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "mu"
	btn.Max = 1.0
	btn.SetLabel("Mu")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "lambda"
	btn.Max = 10.0
	btn.SetLabel("Lambda")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Synapse"
	btn.ev.Field = "alpha"
	btn.Max = 10.0
	btn.SetLabel("Alpha")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	// btnXPos = btnXPos + DefaultButtonWidth + 50
	// wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	// wigBut.SetPos(btnXPos, btnYPos)
	// btn = wigBut.(*ValueButton)
	// btn.ev.Message = "Synapse"
	// btn.ev.Field = "learningRateSlow"
	// btn.Max = 10.0
	// btn.SetLabel("LearnRateSlow")
	// v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	// btn.SetValue(v)
	// iDs++
	// wigBut.SetID(iDs)
	// pw.fields.AddWidget(wigBut)

	// btnXPos = btnXPos + DefaultButtonWidth + 50
	// wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	// wigBut.SetPos(btnXPos, btnYPos)
	// btn = wigBut.(*ValueButton)
	// btn.ev.Message = "Synapse"
	// btn.ev.Field = "learningRateFast"
	// btn.Max = 10.0
	// btn.SetLabel("LearnRateFast")
	// v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	// btn.SetValue(v)
	// iDs++
	// wigBut.SetID(iDs)
	// pw.fields.AddWidget(wigBut)

	iDs++

	return iDs
}

// ############ Neuron row
func (pw *Km0PanelWidget) build_row4(iDs int, btnXPos, btnYPos int32) int {
	wigBut := NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*ValueButton)
	btn.ev.Message = "Neuron"
	btn.ev.Field = "ntao"
	btn.Max = 1000.0
	btn.SetLabel("ntao")
	v := deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Neuron"
	btn.ev.Field = "ntaoS"
	btn.Max = 1000.0
	btn.SetLabel("ntaoS")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Neuron"
	btn.ev.Field = "ntaoJ"
	btn.Max = 1000.0
	btn.SetLabel("ntaoJ")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Neuron"
	btn.ev.Field = "nFastSurge"
	btn.Max = 1000.0
	btn.SetLabel("nFastSurge")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Neuron"
	btn.ev.Field = "nSlowSurge"
	btn.Max = 1000.0
	btn.SetLabel("nSlowSurge")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Neuron"
	btn.ev.Field = "APMax"
	btn.Max = 1000.0
	btn.SetLabel("APMax")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	iDs++

	return iDs
}

// ############ Dendrite row
func (pw *Km0PanelWidget) build_row5(iDs int, btnXPos, btnYPos int32) int {
	wigBut := NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*ValueButton)
	btn.ev.Message = "Dendrite"
	btn.ev.Field = "length"
	btn.Max = 1000.0
	btn.SetLabel("Den length")
	v := deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth + 50
	wigBut = NewValueButton(pw, DefaultValueButtonWidth+50, DefaultValueButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*ValueButton)
	btn.ev.Message = "Dendrite"
	btn.ev.Field = "taoEff"
	btn.Max = 1000.0
	btn.SetLabel("Den taoEff")
	v = deuron.SimModel.GetFloatAsString(btn.ev.Field)
	btn.SetValue(v)
	iDs++
	wigBut.SetID(iDs)
	pw.fields.AddWidget(wigBut)

	iDs++

	return iDs
}
