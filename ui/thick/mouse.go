package thick

import (
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/google/uuid"
)

var mouseMute sync.RWMutex
var onMouseMoveFuncs = map[string]func(pixel.Vec){}

// MouseInput observes mouse activity and calls registered functions.
func MouseInput(win *pixelgl.Window) {
	// Get mouse position if it's in the window
	if win.MouseInsideWindow() {
		prev := win.MousePreviousPosition()
		p := win.MousePosition()

		// If mouse has moved
		if !p.Eq(prev) {
			// Call registered functions
			mouseMute.RLock()
			for _, f := range onMouseMoveFuncs {
				go f(p)
			}
			mouseMute.RUnlock()
		}
	}
}

// OnMouseMove registers a provided function to be called on mouse movement inside the window.
// Returns a function that when called unregisters the OnMouseMove function.
func OnMouseMove(f func(pixel.Vec)) func() {
	mouseMute.Lock()
	id := uuid.New().String()
	onMouseMoveFuncs[id] = f
	mouseMute.Unlock()

	return func() {
		mouseMute.Lock()
		delete(onMouseMoveFuncs, id)
		mouseMute.Unlock()
	}
}
