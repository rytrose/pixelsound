package main

import (
	"image"

	"github.com/faiface/beep"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/nfnt/resize"
)

var pointChan chan image.Point

const traverse = false

func main() {
	pixelgl.Run(run)
}

func run() {
	// Load image
	im, _, err := LoadImageFromFile("images/me.png")

	// Resize image to pretty small
	im = resize.Resize(100, 0, im, resize.NearestNeighbor)
	if err != nil {
		panic(err)
	}

	// Configure UI window
	cfg := pixelgl.WindowConfig{
		Title:  "Pixelsound",
		Bounds: pixel.R(0, 0, float64(im.Bounds().Max.X), float64(im.Bounds().Max.Y)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Create image sprite
	pd := pixel.PictureDataFromImage(im)
	sprite := pixel.NewSprite(pd, pd.Bounds())

	// Create imdraw
	imd := imdraw.New(nil)

	// Create channel for pixel highlighting
	pointChan = make(chan image.Point)

	// Create PixelSound player
	sr := beep.SampleRate(44100)
	player := NewPlayer(win, sr, 2048)

	// Instantiate and play PixelSound
	ps := &PixelSounder{
		T: TtoBLtoR,
		S: SineColor,
	}
	player.SetImagePixelSound(im, ps)

	// PLAY W/TRAVERSAL
	if traverse {
		player.Play(im, ps, 0, 0)
	} else { // PLAY W/MOUSE
		// Do mouse movement
		go OnMouseMove(win, func(p pixel.Vec) {
			player.PlayPixel(p)
		})
	}

	// UI main loop
	for !win.Closed() {
		select {
		case point := <-pointChan:
			DrawImage(win, sprite, imd, im.At(point.X, point.Y), point.X, point.Y)
		default:
		}
		win.Update()
	}
}
