//go:build js

package browser

import (
	"log"

	"github.com/nfnt/resize"
	"github.com/rytrose/pixelsound/ui"
)

type browser struct{}

// Returns a new browser UI for running on the web.
func NewBrowser() ui.UI {
	return &browser{}
}

// Runs the browser UI using canvas.
func (b *browser) Run() {
	// Load image
	imageFilename := "https://i1.sndcdn.com/avatars-000504525639-z2p212-t500x500.jpg"
	im, _, err := LoadImageFromURL(imageFilename)
	if err != nil {
		log.Fatalf("unable to load image %s: %s", imageFilename, err)
	}

	// Resize image to pretty small
	im = resize.Resize(100, 0, im, resize.NearestNeighbor)
	if err != nil {
		log.Fatalf("unable to resize image: %s", err)
	}

	// TODO: create canvas
}
