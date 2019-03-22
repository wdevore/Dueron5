package graphs

import (
	"image"
	"image/color"

	"github.com/wdevore/Deuron5/deuron/app/gui"

	"github.com/fogleman/gg"
	"github.com/veandco/go-sdl2/sdl"
)

// |.   .. .. ...   .  ..  . . ..
// | . . .   .. . . ..  . . . ..
//
//

// PointsGraph renders spikes as pixels.
type PointsGraph struct {
	BaseGraph
	gui.BaseWidget

	dc     *gg.Context
	pixels *image.RGBA

	originX float64
	originY float64

	accessor SeriesAccessor
}

// NewPointsGraph creates a new sdl raw graph
func NewPointsGraph(renderer *sdl.Renderer, texture *sdl.Texture, width, height int) IGraph {
	g := new(PointsGraph)
	g.BaseWidget.Initialize(nil, width, height)
	g.SetGraphics(renderer, texture)
	g.MarkDirty(true)

	// rect := image.Rectangle{image.Point{0, 0}, image.Point{width, height}}
	// g.pixels = image.NewRGBA(rect)
	// g.dc = gg.NewContextForRGBA(g.pixels)

	g.dc = gg.NewContext(width, height)
	g.pixels = g.dc.Image().(*image.RGBA)

	g.originX = 10.0
	g.originY = 10.0

	return g
}

// SetSeries set the chart data
func (g *PointsGraph) SetSeries(accessor SeriesAccessor) {
	g.accessor = accessor

	g.MarkDirty(true)
}

func (g *PointsGraph) GetAccessor() SeriesAccessor {
	ind := 0
	series1XValues := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0}
	series1YValues := []float64{1.0, 20.0 * ran.Float64(), 9.0, 4.0, 2.0, 1.0, 12.0 * ran.Float64()}
	pointColor := color.RGBA{0, 0, 255, 255}

	return func() (x, y float64, c color.Color, more int) {
		series1YValues[1] = 20.0 * ran.Float64()

		if ind == len(series1XValues) {
			ind = 0
			return 0, 0, nil, 0
		}

		x = series1XValues[ind]
		y = series1YValues[ind]
		ind++

		return x, y, pointColor, 1
	}
}

// Destroy release resources
func (g *PointsGraph) Destroy() {
	// err := g.texture.Destroy()
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

// DrawAt renders graph to texture
func (g *PointsGraph) Draw() {
	// g.rect.X = x
	// g.rect.Y = y

	if g.dirty {
		// Render the graph onto the image pixels
		g.dc.SetRGB(1, 1, 1)
		g.dc.Clear()

		g.dc.Identity()
		g.dc.InvertY()

		g.dc.Translate(g.originX, g.originY)

		// Left border
		g.dc.SetRGBA(0.85, 0.85, 0.85, 1.0)
		g.dc.MoveTo(0.0, 0.0)
		g.dc.LineTo(0.0, float64(g.Rect.H)-g.originY*2)
		g.dc.Stroke()

		px, py, c, more := g.accessor()

		for more == 1 {
			g.dc.SetColor(c)
			g.dc.DrawCircle(px, py, 0.75)

			g.dc.Fill()

			px, py, c, more = g.accessor()
		}

		g.texture.Update(&g.Rect, g.pixels.Pix, g.pixels.Stride)

		// Now copy the texture onto the target (aka the display)
		g.renderer.Copy(g.texture, &g.Rect, &g.Rect)

		g.MarkDirty(false)
	}
}

func (g *PointsGraph) Check() bool {
	return false
}
