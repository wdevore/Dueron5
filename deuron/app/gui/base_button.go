package gui

import (
	"image/color"
)

const (
	TextHeight = 10
)

// Toggle background based on state

type baseButton struct {
	BaseWidget

	labelColor  color.RGBA
	borderColor color.RGBA

	// The baseline of the text is at the "bottom" of the text NOT top.
	label string

	textWidth  float64
	textHeight float64
}

func (bb *baseButton) initialize(parent IWidget, width, height int) {
	bb.BaseWidget.Initialize(parent, width, height)

	bb.labelColor = color.RGBA{255, 255, 255, 255}
	bb.borderColor = color.RGBA{255, 255, 255, 255}
}

func (bb *baseButton) SetLabel(label string) {
	bb.label = label
	bb.textWidth, bb.textHeight = bb.DC.MeasureString(label)
}

func (bb *baseButton) isInside(vx, vy int32) bool {
	px, py := bb.parent.Position()
	lx, ly := ViewSpaceToLocal(vx, vy, px, py)

	// fmt.Printf("%d, %d -> %d, %d\n", x, y, lx, ly)
	return PointInside(lx, ly, bb.Rect.X, bb.Rect.Y, bb.Rect.W, bb.Rect.H)

	return true
}
