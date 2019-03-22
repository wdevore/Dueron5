package graphs

import (
	"errors"
	"image"
)

// RGBAWriter is a special type of io.Writer that produces a final image.
type RGBAWriter struct {
	rgba *image.RGBA
}

func NewRGBAWriter() *RGBAWriter {
	w := new(RGBAWriter)
	return w
}

func (ir *RGBAWriter) Write(buffer []byte) (int, error) {
	return 0, nil
}

// SetRGBA sets a raw version of the image.
func (ir *RGBAWriter) SetRGBA(i *image.RGBA) {
	ir.rgba = i
}

// Image returns an *image.Image for the result.
func (ir *RGBAWriter) Image() (*image.RGBA, error) {
	if ir.rgba != nil {
		return ir.rgba, nil
	}
	return nil, errors.New("no valid sources for image data, cannot continue")
}
