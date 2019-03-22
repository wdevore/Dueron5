package gui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

// Layout:
// |  Run  |  RunPause  |  Pause  |  Stop  |  Resume  |  Km0  |  Km1  |  Km2  |  Km3  |
// |  .001 |  .01       |  .1     |  .5    |  1       |  5    |  10   |  25   |  100  |

type PanelWidget struct {
	basePanel

	simControlGroup *ButtonGroup
	keymapGroup     *ButtonGroup
	incDecGroup     *ButtonGroup
	saveGroup       *ButtonGroup

	// Toggle buttons
	autoRunPause IWidget
	rangeSync    IWidget
}

func NewPanelWidget(startID int, renderer *sdl.Renderer, texture *sdl.Texture, width, height int) (widget IWidget, id int) {
	pw := new(PanelWidget)
	pw.initialize(pw, renderer, texture, width, height)
	pw.visible = true
	pw.id = startID
	iDs := pw.build()

	return pw, iDs
}

// override's baseWidget and implements IWidget
func (pw *PanelWidget) Handle(x, y int32, eventType events.MouseEventType) (handled bool, id int) {
	handled, id = pw.simControlGroup.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	handled, _ = pw.keymapGroup.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	handled, _ = pw.incDecGroup.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	handled, _ = pw.saveGroup.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	handled, _ = pw.autoRunPause.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	handled, _ = pw.rangeSync.Handle(x, y, eventType)
	if handled {
		return true, id
	}

	// fmt.Printf("%v, %d\n", handled, id)
	return false, -1
}

func (pw *PanelWidget) Listen(msg *comm.MessageEvent) {
}

func (pw *PanelWidget) DrawAt(x, y int32) {
	pw.SetPos(x, y)
	pw.Draw()
}

func (pw *PanelWidget) Draw() {
	if !pw.visible {
		return
	}

	pw.DC.SetColor(pw.backgroundColor)
	pw.DC.Clear()

	// Draw child widgets
	it := pw.simControlGroup.Iterator()
	for it.Next() {
		btn := it.Value().(*Button)
		btn.Draw()
	}

	it = pw.keymapGroup.Iterator()
	for it.Next() {
		btn := it.Value().(*Button)
		btn.Draw()
	}

	it = pw.incDecGroup.Iterator()
	for it.Next() {
		btn := it.Value().(*Button)
		btn.Draw()
	}

	it = pw.saveGroup.Iterator()
	for it.Next() {
		btn := it.Value().(*Button)
		btn.Draw()
	}

	pw.autoRunPause.Draw()
	pw.rangeSync.Draw()

	pw.update()
}

func (pw *PanelWidget) build() int {
	iDs := pw.id
	btnXPos := int32(0)
	btnYPos := int32(0)

	iDs = pw.build_sim_buttons(iDs, btnXPos, btnYPos)

	// ---------------------------------------------------------------------
	btnXPos = btnXPos + DefaultButtonWidth*7
	iDs = pw.build_panel_buttons(iDs, btnXPos, btnYPos)

	btnXPos = btnXPos + DefaultButtonWidth*6
	iDs = pw.build_saveload_buttons(iDs, btnXPos, btnYPos)

	btnXPos = 0
	btnYPos = DefaultButtonHeight
	iDs = pw.build_model_buttons(iDs, btnXPos, btnYPos)

	return iDs
}

func (pw *PanelWidget) build_sim_buttons(iDs int, btnXPos, btnYPos int32) int {
	pw.simControlGroup = NewButtonGroup()

	wigBut := NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	btn := wigBut.(*Button)
	btn.flashable = true
	btn.ev.Message = "Run"
	btn.SetLabel(btn.ev.Message)
	iDs++
	wigBut.SetID(iDs)
	pw.simControlGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.flashable = true
	btn.ev.Message = "RunPause"
	btn.SetLabel(btn.ev.Message)
	iDs++
	wigBut.SetID(iDs)
	pw.simControlGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.flashable = true
	btn.ev.Message = "Pause"
	btn.SetLabel(btn.ev.Message)
	iDs++
	wigBut.SetID(iDs)
	pw.simControlGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.flashable = true
	btn.ev.Message = "Stop"
	btn.SetLabel(btn.ev.Message)
	iDs++
	wigBut.SetID(iDs)
	pw.simControlGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.flashable = true
	btn.ev.Message = "Resume"
	btn.SetLabel(btn.ev.Message)
	iDs++
	wigBut.SetID(iDs)
	pw.simControlGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	pw.autoRunPause = NewToggleButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	pw.autoRunPause.SetPos(btnXPos, btnYPos)
	togBtn := pw.autoRunPause.(*ToggleButton)
	togBtn.ev.Field = "AutoRunPause"
	selected := deuron.SimModel.GetFloat(togBtn.ev.Field)
	if selected == 1.0 {
		togBtn.Select()
	} else {
		togBtn.UnSelect()
	}
	togBtn.SetLabel("AutoRun")
	iDs++
	wigBut.SetID(iDs)

	btnXPos = btnXPos + DefaultButtonWidth
	pw.rangeSync = NewToggleButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	pw.rangeSync.SetPos(btnXPos, btnYPos)
	togBtn = pw.rangeSync.(*ToggleButton)
	// This command is recognized by the Model listener.
	togBtn.ev.Field = "RangeSync"
	selected = deuron.SimModel.GetFloat(togBtn.ev.Field)
	if selected == 1.0 {
		togBtn.Select()
	} else {
		togBtn.UnSelect()
	}
	togBtn.SetLabel("Range Sync")
	iDs++
	wigBut.SetID(iDs)

	iDs++

	return iDs
}

func (pw *PanelWidget) build_panel_buttons(iDs int, btnXPos, btnYPos int32) int {
	pw.keymapGroup = NewButtonGroup()

	wigBut := NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*Button)
	// This command is recognized by the Gui listener
	btn.ev.Target = "Gui"
	btn.ev.Action = "Hide"
	btn.ev.Message = "Panel"
	btn.SetLabel(btn.ev.Action)
	iDs++
	wigBut.SetID(iDs)
	pw.keymapGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Gui"
	btn.ev.Action = "Show"
	btn.ev.Message = "Panel"
	btn.ev.Value = "Km0"
	btn.Select()
	btn.SetLabel("Main Pan")
	iDs++
	wigBut.SetID(iDs)
	pw.keymapGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Gui"
	btn.ev.Action = "Show"
	btn.ev.Message = "Panel"
	btn.ev.Value = "Km1"
	btn.SetLabel("Misc Pan")
	iDs++
	wigBut.SetID(iDs)
	pw.keymapGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Gui"
	btn.ev.Action = "Show"
	btn.ev.Message = "Panel"
	btn.ev.Value = "Km2"
	btn.SetLabel("Stims Pan")
	iDs++
	wigBut.SetID(iDs)
	pw.keymapGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Gui"
	btn.ev.Action = "Show"
	btn.ev.Message = "Panel"
	btn.ev.Value = "Panel3"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.keymapGroup.AddButton(wigBut)

	iDs++

	return iDs
}

func (pw *PanelWidget) build_saveload_buttons(iDs int, btnXPos, btnYPos int32) int {
	pw.saveGroup = NewButtonGroup()

	wigBut := NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*Button)
	btn.flashable = true
	btn.ev.Target = "App"
	btn.ev.Action = "Command"
	btn.ev.Message = "Save"
	btn.SetLabel("Save")
	iDs++
	wigBut.SetID(iDs)
	pw.saveGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth

	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.flashable = true
	btn.ev.Target = "App"
	btn.ev.Action = "Command"
	btn.ev.Message = "Load"
	btn.SetLabel("Load")
	iDs++
	wigBut.SetID(iDs)
	pw.saveGroup.AddButton(wigBut)

	btnXPos = btnXPos - DefaultButtonWidth
	btnYPos = DefaultButtonHeight
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.flashable = true
	btn.ev.Target = "Graph"
	btn.ev.Action = "Toggle"
	btn.ev.Message = "SpikeGraph"
	btn.ev.Value = "PoissonSamples"
	btn.SetLabel("Poi Sams")
	iDs++
	wigBut.SetID(iDs)
	pw.saveGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth

	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.flashable = true
	btn.ev.Target = "Graph"
	btn.ev.Action = "Toggle"
	btn.ev.Message = "SpikeGraph"
	btn.ev.Value = "StimSamples"
	btn.SetLabel("Stim Sams")
	iDs++
	wigBut.SetID(iDs)
	pw.saveGroup.AddButton(wigBut)

	iDs++

	return iDs
}

func (pw *PanelWidget) build_model_buttons(iDs int, btnXPos, btnYPos int32) int {
	pw.incDecGroup = NewButtonGroup()

	wigBut := NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn := wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "0.001"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "0.01"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "0.1"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "0.5"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.Select()
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "1.0"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "2.0"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "5.0"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "10.0"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "15.0"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "25.0"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	btnXPos = btnXPos + DefaultButtonWidth
	wigBut = NewButton(pw, DefaultButtonWidth, DefaultButtonHeight)
	wigBut.SetPos(btnXPos, btnYPos)
	btn = wigBut.(*Button)
	btn.ev.Target = "Model"
	btn.ev.Action = "Set"
	btn.ev.Field = "Inc/Dec"
	btn.ev.Value = "100.0"
	btn.SetLabel(btn.ev.Value)
	iDs++
	wigBut.SetID(iDs)
	pw.incDecGroup.AddButton(wigBut)

	iDs++

	return iDs
}
