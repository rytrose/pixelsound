//go:build js

package browser

import (
	"bytes"
	"image"
	"math"
	"syscall/js"
	"time"

	"github.com/faiface/beep"
	canvas "github.com/oskca/gopherjs-canvas"
	canvasDOM "github.com/oskca/gopherjs-dom"
	"github.com/rytrose/pixelsound/api"
	"github.com/rytrose/pixelsound/player"
	"github.com/rytrose/pixelsound/sonification"
	"github.com/rytrose/pixelsound/traversal"
	"github.com/rytrose/pixelsound/ui"
	"honnef.co/go/js/dom"
)

type mode int

const (
	modeAlgorithm mode = iota
	modeMouse
	modeKeyboard
)

type bytesReaderCloser struct {
	*bytes.Reader
}

func (r bytesReaderCloser) Close() error { return nil }

type browser struct {
	w        dom.Window
	d        dom.Document
	cv       *canvas.Canvas
	cvc      *canvas.Context2D
	cvEl     dom.Element
	cvIm     *canvasDOM.Element
	imd      image.Image       // Image scaled to be displayed
	imt      image.Image       // Image scaled to be traversed
	r        bytesReaderCloser // Audio file reader
	ext      string            // Audio file extension
	s        api.SonifyFunc
	t        api.TraverseFunc
	mode     mode
	modeChan chan mode
	player   *player.Player
}

// Returns a new browser UI for running on the web.
func NewBrowser() ui.UI {
	// Create DOM references and elements
	w := dom.GetWindow()
	d := w.Document()

	return &browser{
		w:        w,
		d:        d,
		modeChan: make(chan mode),
	}
}

// Run sets the run function on the global javascript scope.
func (b *browser) Run() {
	js.Global().Set("golangRun", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go b.run()
		return nil
	}))
}

// run powers pixelsound on an HTML canvas.
func (b *browser) run() {
	// TODO: refactor
	return

	// Setup player
	sr := beep.SampleRate(44100)
	b.player = player.NewPlayer(sr, 2048, player.WithLatestPoint())

	// Setup canvas elements
	b.cvEl = b.d.GetElementByID("pixelsound")
	b.cv = canvas.New(b.cvEl.Underlying())
	b.cvc = b.cv.GetContext2D()

	// Setup input readers
	go b.readAudioFilesFromInput()
	go b.readImageFilesFromInput()
	go b.readModeFromInput()

	// Kick off animation loop
	b.w.RequestAnimationFrame(b.animate)

	var removeMouseListener func()
	var removeKeyboardListener func()
	for newMode := range b.modeChan {
		// Set current mode
		b.mode = newMode

		// Stop current playback
		b.player.Stop()

		// Clear input listeners
		if removeMouseListener != nil {
			removeMouseListener()
		}
		if removeKeyboardListener != nil {
			removeKeyboardListener()
		}

		// Reset canvas
		b.resetCanvas()

		switch b.mode {
		case modeAlgorithm:
			b.player.Play(b.imt, &api.PixelSounder{
				T: traversal.Random,
				S: sonification.NewAudioScrubber(b.r, b.ext),
			}, image.Point{0, 0})
		case modeMouse:
			lastTraversalPoint := image.Point{0, 0}
			removeMouseListener = OnMouseMove(b.cvEl, func(p image.Point) {
				traversalPoint := image.Point{
					X: int(math.Floor((float64(p.X) / float64(b.imd.Bounds().Dx())) * float64(b.imt.Bounds().Dx()))),
					Y: int(math.Floor((float64(p.Y) / float64(b.imd.Bounds().Dy())) * float64(b.imt.Bounds().Dy()))),
				}
				if (traversalPoint.X != lastTraversalPoint.X) ||
					(traversalPoint.Y != lastTraversalPoint.Y) {
					lastTraversalPoint = traversalPoint
					go b.player.PlayPixel(traversalPoint, false)
				}
			})
		case modeKeyboard:
			lastTraversalPoint := image.Point{0, 0}
			// Write point at origin to update highlight
			b.player.PointLock.Lock()
			b.player.LatestPoint = &lastTraversalPoint
			b.player.PointLock.Unlock()
			removeKeyboardListener = OnKeyboardMove(b.w, func(k keyCode) {
				var p image.Point
				switch k {
				case keyLeft:
					dx := b.imt.Bounds().Dx()
					p = image.Point{
						X: ((((lastTraversalPoint.X - 1) % dx) + dx) % dx),
						Y: lastTraversalPoint.Y,
					}
				case keyRight:
					dx := b.imt.Bounds().Dx()
					p = image.Point{
						X: ((((lastTraversalPoint.X + 1) % dx) + dx) % dx),
						Y: lastTraversalPoint.Y,
					}
				case keyUp:
					dy := b.imt.Bounds().Dy()
					p = image.Point{
						X: lastTraversalPoint.X,
						Y: ((((lastTraversalPoint.Y - 1) % dy) + dy) % dy),
					}
				case keyDown:
					dy := b.imt.Bounds().Dy()
					p = image.Point{
						X: lastTraversalPoint.X,
						Y: ((((lastTraversalPoint.Y + 1) % dy) + dy) % dy),
					}
				}
				lastTraversalPoint = p
				b.player.PlayPixel(p, false)
			})
		}
	}
}

func (b *browser) animate(t time.Duration) {
	b.player.PointLock.HighPriorityLock()
	point := b.player.LatestPoint
	if point != nil {
		b.player.LatestPoint = nil
		b.player.PointLock.HighPriorityUnlock()

		// Update highlight
		b.resetCanvas()
		b.drawPointHighlight(*point)
	} else {
		b.player.PointLock.HighPriorityUnlock()
	}

	// Schedule the next frame
	b.w.RequestAnimationFrame(b.animate)
}

func (b *browser) readModeFromInput() {
	onChange := func(e dom.Event) {
		var newMode mode
		switch e.CurrentTarget().GetAttribute("id") {
		case "algorithm":
			newMode = modeAlgorithm
		case "mouse":
			newMode = modeMouse
		case "keyboard":
			newMode = modeKeyboard
		}
		b.modeChan <- newMode
	}

	for _, id := range []string{"algorithm", "mouse", "keyboard"} {
		el := b.d.GetElementByID(id)
		el.AddEventListener("change", false, onChange)
	}
}

func (b *browser) reloadCurrentMode() {
	b.modeChan <- b.mode
}
