package main

import (
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/google/uuid"
)

var mouseRefreshTimeMs = 16.666666667

var onMouseMoveFuncs = map[string]func(pixel.Vec){}
var onMouseMoveFuncsMutex sync.RWMutex

// MouseInput observes mouse activity and calls registered functions.
func MouseInput(win *pixelgl.Window) {
	t := time.Tick(time.Duration(mouseRefreshTimeMs) * time.Millisecond)

	for {
		// Get mouse position if it's in the window
		if win.MouseInsideWindow() {
			prev := win.MousePreviousPosition()
			p := win.MousePosition()

			// If mouse has moved
			if !p.Eq(prev) {
				// Call registered functions
				onMouseMoveFuncsMutex.RLock()
				for _, f := range onMouseMoveFuncs {
					go f(p)
				}
				onMouseMoveFuncsMutex.RUnlock()
			}
		}

		// Wait refresh time
		<-t
	}
}

// OnMouseMove registers a provided function to be called on mouse movement inside the window.
// Returns a function that when called unregisters the OnMouseMove function.
func OnMouseMove(f func(pixel.Vec)) func() {
	onMouseMoveFuncsMutex.Lock()
	id := uuid.New().String()
	onMouseMoveFuncs[id] = f
	onMouseMoveFuncsMutex.Unlock()

	return func() {
		onMouseMoveFuncsMutex.Lock()
		delete(onMouseMoveFuncs, id)
		onMouseMoveFuncsMutex.Unlock()
	}
}
