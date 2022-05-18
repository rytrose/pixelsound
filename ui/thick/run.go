//go:build !js

package thick

import (
	"flag"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/faiface/beep"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/nfnt/resize"
	"github.com/rytrose/pixelsound/api"
	"github.com/rytrose/pixelsound/log"
	"github.com/rytrose/pixelsound/player"
	"github.com/rytrose/pixelsound/sonification"
	"github.com/rytrose/pixelsound/traversal"
	"github.com/rytrose/pixelsound/ui"
)

type thickClient struct{}

// Returns a new thick client UI for running on Mac OS.
func NewThickClient() ui.UI {
	return &thickClient{}
}

// Runs the Mac OS UI using OpenGL.
func (c *thickClient) Run() {
	pixelgl.Run(c.run)
}

// run implements the client.
func (c *thickClient) run() {
	// Read in command line args
	imageFilename := flag.String("im", "images/me.png", "image to pixelsound")
	inputAudioFilename := flag.String("audio", "audio_inputs/my_name_is_doug_dimmadome.mp3", "audio file to use for pixelsound (if needed)")
	mouse := flag.Bool("mouse", false, "use the mouse to play pixels instead of traverse function")
	keyboard := flag.Bool("keyboard", false, "use the keyboard to play pixels instead of traverse function")
	queue := flag.Bool("queue", false, "all pixels moused over or key pressed to are played sequentially, as opposed to the most recent pixel only")
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

	// Start observing input
	go MouseInput(win)

	// Create image sprite
	pd := pixel.PictureDataFromImage(im)
	sprite := pixel.NewSprite(pd, pd.Bounds())

	// Create imdraw
	imd := imdraw.New(nil)

	// Setup necessary audio inputs
	audioFilename := *inputAudioFilename
	var f *os.File
	var ext string
	if audioFilename != "" {
		f, err = os.Open(audioFilename)
		if err != nil {
			log.Fatalf("unable to open file %s: %s", audioFilename, err)
		}
		split := strings.Split(audioFilename, ".")
		ext = strings.ToLower(split[len(split)-1])
	}

	// Create PixelSound player
	sr := beep.SampleRate(44100)
	player := player.NewPlayer(sr, 2048)

	// Instantiate and play PixelSound
	t, ok := traversal.TraverseFuncs[*traverseFunc]
	if !ok {
		log.Fatalf("no traversal function named %s", *traverseFunc)
	}
	_, ok = sonification.SonifyFuncNames[*sonifyFunc]
	if !ok {
		log.Fatalf("no sonification function named %s", *sonifyFunc)
	}
	var s api.SonifyFunc
	switch *sonifyFunc {
	case "SineColor":
		s = sonification.NewSineColor(sr)
	case "AudioScrubber":
		s = sonification.NewAudioScrubber(f, ext)
	}
	ps := &api.PixelSounder{
		T: t,
		S: s,
	}
	player.SetImagePixelSound(im, ps)

	// PLAY W/MOUSE
	if *mouse {
		// Register play pixel on mouse movement
		stop := OnMouseMove(func(p pixel.Vec) {
			player.PlayPixel(convertPixelToImage(win, p), *queue)
		})
		defer stop()
	} else if *keyboard {
		// Setup play pixel by arrow keys
		keyboardPixelLocation := pixel.Vec{
			X: 0,
			Y: win.Bounds().Max.Y,
		}

		// LEFT ARROW
		stopL := OnKeyPress(pixelgl.KeyLeft, func(b pixelgl.Button) {
			newX := keyboardPixelLocation.X - 1
			if newX < 0 {
				newX = win.Bounds().Max.X
			}
			keyboardPixelLocation.X = newX
			player.PlayPixel(convertPixelToImage(win, keyboardPixelLocation), *queue)
		}, true)
		defer stopL()

		// RIGHT ARROW
		stopR := OnKeyPress(pixelgl.KeyRight, func(b pixelgl.Button) {
			newX := keyboardPixelLocation.X + 1
			if newX > win.Bounds().Max.X {
				newX = 0
			}
			keyboardPixelLocation.X = newX
			player.PlayPixel(convertPixelToImage(win, keyboardPixelLocation), *queue)
		}, true)
		defer stopR()

		// UP ARROW
		stopU := OnKeyPress(pixelgl.KeyUp, func(b pixelgl.Button) {
			newY := keyboardPixelLocation.Y + 1
			if newY > win.Bounds().Max.Y {
				newY = 0
			}
			keyboardPixelLocation.Y = newY
			player.PlayPixel(convertPixelToImage(win, keyboardPixelLocation), *queue)
		}, true)
		defer stopU()

		// DOWN ARROW
		stopD := OnKeyPress(pixelgl.KeyDown, func(b pixelgl.Button) {
			newY := keyboardPixelLocation.Y - 1
			if newY < 0 {
				newY = win.Bounds().Max.Y
			}
			keyboardPixelLocation.Y = newY
			player.PlayPixel(convertPixelToImage(win, keyboardPixelLocation), *queue)
		}, true)
		defer stopD()
	} else { // PLAY W/TRAVERSAL
		player.Play(im, ps, image.Point{0, 0})
	}

	// Draw initial picture
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	// UI main loop
	for !win.Closed() {
		select {
		case point := <-player.PointChan:
			DrawImage(win, sprite, imd, im.At(point.X, point.Y), point)
		default:
		}
		win.Update()
		MouseInput(win)
		KeyboardUpdate(win)
	}
}
