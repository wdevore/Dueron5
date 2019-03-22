package gui

import (
	"fmt"

	"github.com/fogleman/gg"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/wdevore/Deuron5/deuron/app/comm"
)

const (
	DefaultButtonWidth       = 75
	DefaultButtonHeight      = 25
	DefaultValueButtonWidth  = 75
	DefaultValueButtonHeight = 50
)

type BaseWidget struct {
	// Message bus fields
	ev comm.MessageEvent

	id int
	// The parent is for children.
	parent IWidget

	visible bool
	canHide bool

	Rect sdl.Rect

	DC *gg.Context

	// position
	x, y int32
}

func (bw *BaseWidget) Initialize(parent IWidget, width, height int) {

	bw.parent = parent
	bw.Rect.W = int32(width)
	bw.Rect.H = int32(height)

	if parent != nil {
		if parent.Context() == nil {
			bw.DC = gg.NewContext(width, height)
		} else {
			bw.DC = parent.Context()
		}
	}
}

func (bw *BaseWidget) Handle(x, y int32) (handled bool, id int) {
	return false, -1
}

func (bw *BaseWidget) Refresh() {
}

func (bw *BaseWidget) DrawAt(x, y int32) {
	bw.Rect.X = x
	bw.Rect.Y = y
}

func (bw *BaseWidget) Draw() {
}

func (bw *BaseWidget) Show() {
	bw.visible = true
}

func (bw *BaseWidget) Hide() {
	if bw.canHide {
		bw.visible = false
	}
}

func (bw *BaseWidget) IsVisible() bool {
	return bw.visible
}

func (bw *BaseWidget) SetID(id int) {
	bw.id = id
	bw.ev.ID = fmt.Sprintf("%d", id)
}

func (bw *BaseWidget) ID() int {
	return bw.id
}

func (bw *BaseWidget) Position() (x, y int32) {
	return bw.x, bw.y
}

func (bw *BaseWidget) SetPos(x, y int32) {
	bw.Rect.X = x
	bw.Rect.Y = y
	bw.x = x
	bw.y = y
}

func (bw *BaseWidget) SetContext(dc *gg.Context) {
	bw.DC = dc
}

func (bw *BaseWidget) Context() *gg.Context {
	return bw.DC
}
