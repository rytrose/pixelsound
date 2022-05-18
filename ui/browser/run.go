//go:build js

package browser

import (
	"fmt"
	"image"
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

type browser struct {
	w      dom.Window
	d      dom.Document
	cv     *canvas.Canvas
	cvc    *canvas.Context2D
	cvIm   *canvasDOM.Element
	im     image.Image
	player *player.Player
}

// Returns a new browser UI for running on the web.
func NewBrowser() ui.UI {
	w := dom.GetWindow()
	d := w.Document()
	cvHTML := d.GetElementByID("pixelsound")
	cv := canvas.New(cvHTML.Underlying())
	cv.Width = 100
	cv.Height = 100
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

	// Resize image to pretty small
	b.im = resize.Resize(100, 0, im, resize.NearestNeighbor)
	if err != nil {
		log.Fatalf("unable to resize image: %s", err)
	}

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
	b.player.SetImagePixelSound(b.im, ps)
	b.player.Play(b.im, ps, image.Point{0, 0})

	// Draw initial image
	b.drawImage()

	// Kick off animation loop
	b.w.RequestAnimationFrame(b.run)
}

func (b *browser) run(t time.Duration) {
	select {
	case point := <-b.player.PointChan:
		b.cvc.ClearRect(0, 0, 100, 100)
		b.drawImage()
		b.highlightPoint(point)
	default:
	}

	// Schedule the next frame
	b.w.RequestAnimationFrame(b.run)
}

func (b *browser) drawImage() {
	b.cvc.DrawImage(b.cvIm, 0, 0, 100, 100)
}

func (b *browser) highlightPoint(p image.Point) {
	b.cvc.StrokeStyle = "#000"
	b.cvc.LineWidth = 1.0
	b.cvc.StrokeRect(float64(p.X-6), float64(p.Y-6), 12, 12)
	red, green, blue, _ := util.Uint8RGBA(b.im.At(p.X, p.Y))
	b.cvc.FillStyle = fmt.Sprintf("rgb(%d, %d, %d)", red, green, blue)
	b.cvc.FillRect(float64(p.X-5), float64(p.Y-5), 10, 10)
}
