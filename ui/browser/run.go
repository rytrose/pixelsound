//go:build js

package browser

import (
	"bytes"
	"image"
	"math"
	"sync/atomic"
	"syscall/js"
	"time"

	"github.com/faiface/beep"
	canvas "github.com/oskca/gopherjs-canvas"
	canvasDOM "github.com/oskca/gopherjs-dom"
	"github.com/rytrose/pixelsound/api"
	"github.com/rytrose/pixelsound/log"
	"github.com/rytrose/pixelsound/player"
	"github.com/rytrose/pixelsound/sonification"
	"github.com/rytrose/pixelsound/traversal"
	"github.com/rytrose/pixelsound/ui"
	"honnef.co/go/js/dom"
)

type loadingState int32

const (
	notLoading loadingState = iota
	loading
)

type bytesReaderCloser struct {
	*bytes.Reader
}

func (r bytesReaderCloser) Close() error { return nil }

type browser struct {
	w                   dom.Window
	d                   dom.Document
	cv                  *canvas.Canvas
	cvc                 *canvas.Context2D
	cvEl                dom.Element
	cvIm                *canvasDOM.Element
	im                  image.Image
	r                   *bytesReaderCloser // Audio file reader
	ext                 string             // Audio file extension
	loadingState        loadingState
	removeMouseListener func()
	player              *player.Player
}

// Returns a new browser UI for running on the web.
func NewBrowser() ui.UI {
	// Create DOM references and elements
	w := dom.GetWindow()
	d := w.Document()

	return &browser{
		w:            w,
		d:            d,
		loadingState: notLoading,
	}
}

// Run sets functions on the JS window.
func (b *browser) Run() {
	js.Global().Set("golangSetup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go b.setup()
		return nil
	}))

	js.Global().Set("golangRun", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go b.run()
		return nil
	}))

	js.Global().Set("golangUpdateImage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go b.updateImage(args[0].String())
		return nil
	}))

	js.Global().Set("golangUpdateAudio", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go b.updateAudio(args[0].String())
		return nil
	}))

	js.Global().Call("jsGolangReady")
}

func (b *browser) setup() {
	// Setup player
	sr := beep.SampleRate(44100)
	b.player = player.NewPlayer(sr, 2048, player.WithLatestPoint())
	js.Global().Call("jsGolangSetup")
}

// run powers pixelsound on an HTML canvas.
func (b *browser) run() {
	// Setup canvas elements
	b.cvEl = b.d.GetElementByID("pixelsound")
	b.cv = canvas.New(b.cvEl.Underlying())
	b.cvc = b.cv.GetContext2D()

	// Kick off animation loop
	b.w.RequestAnimationFrame(b.animate)

	// Start mouse-based audio
	updateWaveform := func(progress float64) {
		js.Global().Call("jsUpdateWaveform", progress)
	}
	lastTraversalPoint := image.Point{0, 0}
	b.removeMouseListener = OnMouseMove(b.cvEl, func(p image.Point, width int, height int) {
		// TODO add fidelity slider to "lower the resolution"

		if b.getLoadingState() != loading && b.im != nil && b.r != nil {
			// p is the location relative to the size of the canvas.
			// width and height are the current size of the canvas.
			// Translate the relative location to the corresponding
			// location on the original sized image.
			traversalPoint := image.Point{
				X: int(math.Floor((float64(p.X) / float64(width)) * float64(b.im.Bounds().Dx()))),
				Y: int(math.Floor((float64(p.Y) / float64(height)) * float64(b.im.Bounds().Dy()))),
			}
			if (traversalPoint.X != lastTraversalPoint.X) ||
				(traversalPoint.Y != lastTraversalPoint.Y) {
				lastTraversalPoint = traversalPoint
				go b.player.PlayPixel(traversalPoint, false, updateWaveform)
			}
		}
	})
}

func (b *browser) animate(t time.Duration) {
	func() {
		b.player.PointLock.HighPriorityLock()
		defer b.player.PointLock.HighPriorityUnlock()
		point := b.player.LatestPoint
		if point != nil {
			// Update highlight
			b.resetCanvas()
			b.drawPointHighlight(*point)
		}
	}()

	// Schedule the next frame
	b.w.RequestAnimationFrame(b.animate)
}

func (b *browser) setLoadingState(newState loadingState) {
	atomic.StoreInt32((*int32)(&b.loadingState), int32(newState))
}

func (b *browser) getLoadingState() loadingState {
	return loadingState(atomic.LoadInt32((*int32)(&b.loadingState)))
}

func (b *browser) updateImage(dataURLString string) {
	defer func() {
		js.Global().Call("jsImageUpdated")
		b.setLoadingState(notLoading)
	}()
	b.setLoadingState(loading)
	im, err := decodeImageFromDataURL(dataURLString)
	if err != nil {
		log.Println("unable to decode image from data URL", err)
		return
	}
	b.im = im
	b.player.SetImage(b.im)
}

func (b *browser) updateAudio(dataURLString string) {
	defer func() {
		js.Global().Call("jsAudioUpdated")
		b.setLoadingState(notLoading)
	}()
	b.setLoadingState(loading)
	data, ext, err := decodeAudioFromDataURL(dataURLString)
	if err != nil {
		log.Println("unable to decode audio from data URL", err)
	}

	// Stop playing audio
	b.player.Stop()

	b.r = &bytesReaderCloser{bytes.NewReader(data)}
	b.ext = ext

	b.player.SetPixelSound(&api.PixelSounder{
		T: traversal.Random,
		S: sonification.NewAudioScrubber(b.r, b.ext),
	})
}
