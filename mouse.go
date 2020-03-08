package main

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var refreshTimeMs = 16.666666667

// OnMouseMove calls a provided function on mouse movement inside the window.
func OnMouseMove(win *pixelgl.Window, f func(pixel.Vec)) {
	t := time.Tick(time.Duration(refreshTimeMs) * time.Millisecond)

	for {
		// Get mouse position if it's in the window
		if win.MouseInsideWindow() {
			prev := win.MousePreviousPosition()
			p := win.MousePosition()

			// If mouse has moved
			if !p.Eq(prev) {
				// Call function
				go f(p)
			}
		}

		// Wait refresh time
		<-t
	}
}
