package gui

import (
	"github.com/fogleman/gg"
	"github.com/wdevore/Deuron5/deuron/app/events"
)

type IWidget interface {
	SetID(int)
	ID() int

	Handle(x, y int32, eventType events.MouseEventType) (bool, int)
	Draw()
	DrawAt(x, y int32)

	Refresh()

	Position() (x, y int32)
	SetPos(x, y int32)

	SetContext(*gg.Context)
	Context() *gg.Context

	Show()
	Hide()
	IsVisible() bool
}

// Map view-space coords (aka px, py) to local-space
func ViewSpaceToLocal(vx, vy, px, py int32) (lx, ly int32) {
	// Map view-space coords (aka px, py) to local-space
	lx = vx - px
	ly = vy - py

	return lx, ly
}

func PointInside(x, y, rx, ry, w, h int32) bool {
	if x < rx {
		return false
	}
	if y < ry {
		return false
	}
	if x >= rx+w {
		return false
	}
	if y >= ry+h {
		return false
	}

	return true
}
