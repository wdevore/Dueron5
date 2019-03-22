package gui

import (
	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron/app/comm"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

// Contains panels and graphs
type Gui struct {
	mainPanel IWidget
	km0Panel  IWidget
	km1Panel  IWidget
	km2Panel  IWidget

	visiblePanel IWidget

	panels *sll.List
}

func NewGui(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) *Gui {
	g := new(Gui)

	g.panels = sll.New()

	controlIDs := 0
	g.mainPanel, controlIDs = NewPanelWidget(controlIDs, renderer, texture, 1500, 50)
	g.mainPanel.SetPos(500, 0)
	g.panels.Add(g.mainPanel)

	g.km0Panel, controlIDs = NewKm0PanelWidget(controlIDs+1, renderer, texture, 1000, 500)
	g.km0Panel.SetPos(1000, 55)
	g.km0Panel.Show()
	g.visiblePanel = g.km0Panel
	g.panels.Add(g.km0Panel)

	g.km1Panel, controlIDs = NewKm1PanelWidget(controlIDs+1, renderer, texture, 1000, 500)
	g.km1Panel.SetPos(1000, 55)
	g.panels.Add(g.km1Panel)

	g.km2Panel, controlIDs = NewKm2PanelWidget(controlIDs+1, renderer, texture, 1000, 500)
	g.km2Panel.SetPos(1000, 55)
	g.panels.Add(g.km2Panel)

	// fmt.Printf("IDs generated: %d\n", controlIDs+1)
	return g
}

func (g *Gui) Handle(vx, vy int32, eventType events.MouseEventType) bool {
	it := g.panels.Iterator()
	for it.Next() {
		p := it.Value().(IWidget)
		if p.IsVisible() {
			handled, _ := p.Handle(vx, vy, eventType)
			if handled {
				return true
			}
		}
	}

	return false
}

func (g *Gui) SubscribeToMsgBus() {
	comm.MsgBus.Subscribe(g)

	it := g.panels.Iterator()
	for it.Next() {
		p := it.Value().(comm.IMessageListener)
		comm.MsgBus.Subscribe(p)
	}
}

func (g *Gui) DrawAt(x, y int32) {
	it := g.panels.Iterator()
	for it.Next() {
		p := it.Value().(IWidget)
		p.DrawAt(x, y)
	}
}

func (g *Gui) Draw() {
	it := g.panels.Iterator()
	for it.Next() {
		p := it.Value().(IWidget)
		p.Draw()
	}
}

func (g *Gui) Refresh() {
	it := g.panels.Iterator()
	for it.Next() {
		p := it.Value().(IWidget)
		p.Refresh()
	}
}

func (g *Gui) Listen(msg *comm.MessageEvent) {
	// fmt.Printf("Gui.Listen: %s\n", msg)

	switch msg.Target {
	case "Gui":
		switch msg.Action {
		case "Show":
			it := g.panels.Iterator()
			for it.Next() {
				p := it.Value().(IWidget)
				p.Hide()
			}

			switch msg.Value {
			case "Km0":
				g.visiblePanel = g.km0Panel
				g.visiblePanel.(IWidget).Show()
				return
			case "Km1":
				g.visiblePanel = g.km1Panel
				g.visiblePanel.(IWidget).Show()
				return
			case "Km2":
				g.visiblePanel = g.km2Panel
				g.visiblePanel.(IWidget).Show()
				return
			}
			// fmt.Printf("%s = %s\n", msg.Message, msg.Value)
			// m.SetAsFloat(msg.Message, msg.Value)
			return
		case "Hide":
			g.visiblePanel.Hide()
			return
		}
	}
}
