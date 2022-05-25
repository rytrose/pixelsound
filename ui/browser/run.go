//go:build js

package browser

import (
	"fmt"
	"image"
	"math"
	"time"

	"github.com/faiface/beep"
	"github.com/nfnt/resize"
	canvas "github.com/oskca/gopherjs-canvas"
	canvasDOM "github.com/oskca/gopherjs-dom"
	"github.com/rytrose/pixelsound/api"
	"github.com/rytrose/pixelsound/log"
	"github.com/rytrose/pixelsound/player"
	"github.com/rytrose/pixelsound/sonification"
	"github.com/rytrose/pixelsound/traversal"
	"github.com/rytrose/pixelsound/ui"
	"github.com/rytrose/pixelsound/util"
	"honnef.co/go/js/dom"
)

// Scale of traversal relative to displayed image. Must be between 0 and 1.
// A value of 1 means traversal will be per pixel for the displayed image.
// Values less than 1 traverse the image in fewer steps.
const scale = 0.2
const displayMaxHeight = 200
const displayMaxWidth = 200

type browser struct {
	w      dom.Window
	d      dom.Document
	cv     *canvas.Canvas
	cvc    *canvas.Context2D
	cvIm   *canvasDOM.Element
	imd    image.Image // Image scaled to be displayed
	imt    image.Image // Image scaled to be traversed
	player *player.Player
}

// Returns a new browser UI for running on the web.
func NewBrowser() ui.UI {
	w := dom.GetWindow()
	d := w.Document()
	cvHTML := d.GetElementByID("pixelsound")
	cv := canvas.New(cvHTML.Underlying())
	return &browser{
		w:   w,
		d:   d,
		cv:  cv,
		cvc: cv.GetContext2D(),
	}
}

// Runs the browser UI using canvas.
func (b *browser) Run() {
	// Load image
	imageFilename := "https://i1.sndcdn.com/avatars-000504525639-z2p212-t500x500.jpg"
	im, dataURL, err := LoadImageFromURL(imageFilename)
	if err != nil {
		log.Fatalf("unable to load image %s: %s", imageFilename, err)
	}

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

	// Create PixelSound player
	sr := beep.SampleRate(44100)
	b.player = player.NewPlayer(sr, 2048)
	ps := &api.PixelSounder{
		T: traversal.TtoBLtoR,
		S: sonification.NewSineColor(sr),
	}
	b.player.SetImagePixelSound(b.imt, ps)
	b.player.Play(b.imt, ps, image.Point{0, 0})

	// Draw initial image
	b.drawImage()

	// Kick off animation loop
	b.w.RequestAnimationFrame(b.run)
}

func (b *browser) run(t time.Duration) {
	select {
	case point := <-b.player.PointChan:
		b.cvc.ClearRect(0, 0, float64(b.cv.Width), float64(b.cv.Height))
		b.drawImage()
		b.highlightPoint(point)
	default:
	}

	// Schedule the next frame
	b.w.RequestAnimationFrame(b.run)
}

func (b *browser) drawImage() {
	b.cvc.DrawImage(b.cvIm, 0, 0, float64(b.cv.Width), float64(b.cv.Height))
}

func (b *browser) highlightPoint(p image.Point) {
	red, green, blue, _ := util.Uint8RGBA(b.imt.At(p.X, p.Y))
	displayPoint := image.Point{
		X: int(math.Floor((float64(p.X) / float64(b.imt.Bounds().Dx())) * float64(b.imd.Bounds().Dx()))),
		Y: int(math.Floor((float64(p.Y) / float64(b.imt.Bounds().Dy())) * float64(b.imd.Bounds().Dy()))),
	}

	b.cvc.StrokeStyle = "#000"
	b.cvc.LineWidth = 1.0
	b.cvc.StrokeRect(float64(displayPoint.X-6), float64(displayPoint.Y-6), 12, 12)
	b.cvc.FillStyle = fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue)
	b.cvc.FillRect(float64(displayPoint.X-5), float64(displayPoint.Y-5), 10, 10)
}
