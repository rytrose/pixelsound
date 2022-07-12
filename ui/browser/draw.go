//go:build js

package browser

import (
	"fmt"
	"image"
	"math"

	"github.com/rytrose/pixelsound/util"
)

// Scale of traversal relative to displayed image. Must be between 0 and 1.
// A value of 1 means traversal will be per pixel for the displayed image.
// Values less than 1 traverse the image in fewer steps.
const scale = 0.1

func (b *browser) resetCanvas() {
	b.cvc.ClearRect(0, 0, float64(b.cv.Width), float64(b.cv.Height))
}

// drawPointHighlight takes a point on the canvas and highlights it.
func (b *browser) drawPointHighlight(p image.Point) {
	width := b.cvEl.Underlying().Get("width").Int()
	height := b.cvEl.Underlying().Get("height").Int()

	displayPoint := image.Point{
		X: int(math.Floor((float64(p.X) / float64(b.im.Bounds().Dx())) * float64(width))),
		Y: int(math.Floor((float64(p.Y) / float64(b.im.Bounds().Dy())) * float64(height))),
	}

	b.cvc.StrokeStyle = "#000"
	b.cvc.LineWidth = 1.0
	b.cvc.StrokeRect(float64(displayPoint.X-6), float64(displayPoint.Y-6), 12, 12)

	red, green, blue, _ := util.Uint8RGBA(b.im.At(p.X, p.Y))
	b.cvc.FillStyle = fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue)
	b.cvc.FillRect(float64(displayPoint.X-5), float64(displayPoint.Y-5), 10, 10)
}
