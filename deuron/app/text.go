package app

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

// Text represents a text texture for rendering.
type Text struct {
	nFont    *Font
	renderer *sdl.Renderer

	texture *sdl.Texture
	color   sdl.Color
	text    string
	bounds  sdl.Rect

	setTextMutex *sync.Mutex
}

// NewText creates a Text object.
func NewText(font *Font, renderer *sdl.Renderer) *Text {
	t := new(Text)
	t.nFont = font
	t.renderer = renderer
	t.initialize()

	t.color = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	t.setTextMutex = &sync.Mutex{}
	t.bounds = sdl.Rect{}

	return t
}

// Initialize sets up Text based on TextPath
func (t *Text) initialize() error {
	return nil
}

func (t *Text) Text() string {
	return t.text
}

// SetText builds an SDL texture. Be sure to call Destroy before
// program exit.
func (t *Text) SetText(text string) (err error) {
	t.setTextMutex.Lock()
	defer t.setTextMutex.Unlock()

	if text == "" {
		fmt.Println("Blank text---------------")
		return nil
	}

	t.text = text

	t.texture.Destroy()
	t.texture = nil

	var surface *sdl.Surface

	// First we draw an image to a surface
	surface, err = t.nFont.font.RenderUTF8Solid(text, t.color)
	if err != nil {
		fmt.Printf("surface err: %c, Text `%s`\n", err, text)
		return err
	}

	// Now generate a texture for rendering, using the surface.
	t.texture, err = t.renderer.CreateTextureFromSurface(surface)

	if err != nil {
		fmt.Printf("Tex create err: %c\n", err)
		t.Destroy()
		return err
	}

	_, _, width, height, qerr := t.texture.Query()
	if qerr != nil {
		t.Destroy()
		return qerr
	}

	t.bounds = sdl.Rect{X: 0, Y: 0, W: width, H: height}

	// We don't need the surface any longer now that the texture
	// is created.
	surface.Free()

	return nil
}

// Draw renders text
func (t *Text) Draw() {
	t.setTextMutex.Lock()
	defer t.setTextMutex.Unlock()
	err := t.renderer.Copy(t.texture, nil, &t.bounds)
	if err != nil {
		fmt.Printf("Text::Draw failed copy: (%v), %v\n", err, t.texture)
	}
}

// DrawAt renders text
func (t *Text) DrawAt(x, y int32) {
	t.setTextMutex.Lock()
	defer t.setTextMutex.Unlock()
	t.bounds.X = x
	t.bounds.Y = y
	t.renderer.Copy(t.texture, nil, &t.bounds)
}

// Destroy closes the Text
func (t *Text) Destroy() error {
	t.setTextMutex.Lock()
	defer t.setTextMutex.Unlock()
	if t.texture != nil {
		return t.texture.Destroy()
	}
	return nil
}
