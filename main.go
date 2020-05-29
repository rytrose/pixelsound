package main

import (
	"flag"
	"fmt"
	"image"

	"github.com/faiface/beep"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/nfnt/resize"
)

var pointChan chan image.Point

func main() {
	pixelgl.Run(run)
}

func run() {
	// Read in command line args
	imageFilename := flag.String("im", "images/me.png", "image to pixelsound")
	inputAudioFilename := flag.String("audio", "audio_inputs/my_name_is_doug_dimmadome.mp3", "audio file to use for pixelsound (if needed)")
	mouse := flag.Bool("mouse", false, "use the mouse to play pixels instead of traverse function")
	mouseQueue := flag.Bool("mouseQueue", false, "all pixels moused over are played sequentially, as opposed to the most recent pixel only")
	traverseFunc := flag.String("t", "TtoBLtoR", "traversal function to use")
	sonifyFunc := flag.String("s", "SineColor", "sonification function to use")
	flag.Parse()

	// Load image
	im, _, err := LoadImageFromFile(*imageFilename)
	if err != nil {
		panic(fmt.Sprintf("unable to load image %s: %s", *imageFilename, err))
	}

	// Resize image to pretty small
	im = resize.Resize(100, 0, im, resize.NearestNeighbor)
	if err != nil {
		panic(fmt.Sprintf("unable to resize image: %s", err))
	}

	// Configure UI window
	cfg := pixelgl.WindowConfig{
		Title:  "Pixelsound",
		Bounds: pixel.R(0, 0, float64(im.Bounds().Max.X), float64(im.Bounds().Max.Y)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(fmt.Sprintf("unable to create window: %s", err))
	}

	// Create image sprite
	pd := pixel.PictureDataFromImage(im)
	sprite := pixel.NewSprite(pd, pd.Bounds())

	// Create imdraw
	imd := imdraw.New(nil)

	// Create channel for pixel highlighting
	pointChan = make(chan image.Point)

	// Setup necessary audio inputs
	audioFilename = *inputAudioFilename

	// Create PixelSound player
	sr := beep.SampleRate(44100)
	player := NewPlayer(win, sr, 2048)

	// Instantiate and play PixelSound
	t, ok := traverseFuncs[*traverseFunc]
	if !ok {
		panic(fmt.Sprintf("no traversal function named %s", *traverseFunc))
	}
	s, ok := sonifyFuncs[*sonifyFunc]
	if !ok {
		panic(fmt.Sprintf("no sonification function named %s", *sonifyFunc))
	}
	ps := &PixelSounder{
		T: t,
		S: s,
	}
	player.SetImagePixelSound(im, ps)

	// PLAY W/MOUSE
	if *mouse {
		// Do mouse movement
		go OnMouseMove(win, func(p pixel.Vec) {
			player.PlayPixel(p, *mouseQueue)
		})
	} else { // PLAY W/TRAVERSAL
		player.Play(im, ps, 0, 0)
	}

	// Draw initial picture
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

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
