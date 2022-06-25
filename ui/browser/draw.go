//go:build js

package browser

import (
	"fmt"
	"image"
	"math"

	"github.com/nfnt/resize"
	canvasDOM "github.com/oskca/gopherjs-dom"
	"github.com/rytrose/pixelsound/util"
	"github.com/vincent-petithory/dataurl"
	"honnef.co/go/js/dom"
)

// Scale of traversal relative to displayed image. Must be between 0 and 1.
// A value of 1 means traversal will be per pixel for the displayed image.
// Values less than 1 traverse the image in fewer steps.
const scale = 0.1
const displayMaxHeight = 200
const displayMaxWidth = 200

func (b *browser) resetCanvas() {
	b.cvc.ClearRect(0, 0, float64(b.cv.Width), float64(b.cv.Height))
	b.drawImage()
}

func (b *browser) drawImage() {
	b.cvc.DrawImage(b.cvIm, 0, 0, float64(b.cv.Width), float64(b.cv.Height))
}

// drawPointHighlight takes a point in the traversal scale and highlights it in the display scale.
func (b *browser) drawPointHighlight(p image.Point) {
	displayPoint := image.Point{
		X: int(math.Floor((float64(p.X) / float64(b.imt.Bounds().Dx())) * float64(b.imd.Bounds().Dx()))),
		Y: int(math.Floor((float64(p.Y) / float64(b.imt.Bounds().Dy())) * float64(b.imd.Bounds().Dy()))),
	}

	b.cvc.StrokeStyle = "#000"
	b.cvc.LineWidth = 1.0
	b.cvc.StrokeRect(float64(displayPoint.X-6), float64(displayPoint.Y-6), 12, 12)

	red, green, blue, _ := util.Uint8RGBA(b.imt.At(p.X, p.Y))
	b.cvc.FillStyle = fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue)
	b.cvc.FillRect(float64(displayPoint.X-5), float64(displayPoint.Y-5), 10, 10)
}

func (b *browser) setImage(im image.Image, dataURL *dataurl.DataURL) {
	// Resize image for display
	maxWidthDelta := math.Abs(float64(im.Bounds().Dx()) - float64(displayMaxWidth))
	maxHeightDelta := math.Abs(float64(im.Bounds().Dy()) - float64(displayMaxHeight))
	if maxWidthDelta < maxHeightDelta {
		b.imd = resize.Resize(displayMaxWidth, 0, im, resize.NearestNeighbor)
	} else {
		b.imd = resize.Resize(0, displayMaxHeight, im, resize.NearestNeighbor)
	}
	// Update canvas size based on image
	b.cv.Width = b.imd.Bounds().Dx()
	b.cv.Height = b.imd.Bounds().Dy()

	// Resize image for traversal
	newWidth := uint(math.Floor(scale * float64(im.Bounds().Dx())))
	b.imt = resize.Resize(newWidth, 0, im, resize.NearestNeighbor)

	// Load image
	loadedChan := make(chan struct{})
	imgEl := b.d.CreateElement("img")
	imgEl.SetAttribute("src", dataURL.String())
	imgEl.AddEventListener("load", false, func(el dom.Event) {
		imgElCanvas := canvasDOM.WrapElement(el.Target().Underlying())
		b.cvIm = imgElCanvas
		loadedChan <- struct{}{}
	})
	<-loadedChan
	b.cvc.ClearRect(0, 0, float64(b.cv.Width), float64(b.cv.Height))
	b.drawImage()
}
