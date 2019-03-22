package gui

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/veandco/go-sdl2/sdl"
)

type basePanel struct {
	BaseWidget

	renderer *sdl.Renderer
	texture  *sdl.Texture

	pixels *image.RGBA
	pdc    *gg.Context

	backgroundColor color.RGBA
}

func (bp *basePanel) initialize(parent IWidget, renderer *sdl.Renderer, texture *sdl.Texture, width, height int) {
	// Create Context for all widgets of this panel (aka parent).
	bp.pdc = gg.NewContext(width, height)
	bp.SetContext(bp.pdc)

	bp.BaseWidget.Initialize(parent, width, height)

	bp.renderer = renderer
	bp.texture = texture

	bp.pixels = bp.pdc.Image().(*image.RGBA)

	bp.backgroundColor = color.RGBA{200, 127, 127, 32}

	bp.texture.SetBlendMode(sdl.BLENDMODE_BLEND)
}

func (bp *basePanel) update() {
	bp.texture.Update(&bp.Rect, bp.pixels.Pix, bp.pixels.Stride)

	// Now copy the texture onto the target (aka the display)
	bp.renderer.Copy(bp.texture, &bp.Rect, &bp.Rect)
}
