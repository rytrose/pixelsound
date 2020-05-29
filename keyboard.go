package main

import (
	"sync"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/google/uuid"
)

var keyboardRefreshTimeMs = 16.666666667

var onKeyPressedFuncs = map[pixelgl.Button]map[string]func(pixelgl.Button){}
var onKeyPressedFuncsMutex sync.RWMutex

// KeyboardInput observes keyboard activity and calls registered functions.
func KeyboardInput(win *pixelgl.Window) {
	t := time.Tick(time.Duration(keyboardRefreshTimeMs) * time.Millisecond)

	for {
		// Check if registered keys were pressed
		onKeyPressedFuncsMutex.RLock()
		for key := range onKeyPressedFuncs {
			if win.JustPressed(key) {
				for _, f := range onKeyPressedFuncs[key] {
					go f(key)
				}
			}
		}
		onKeyPressedFuncsMutex.RUnlock()

		// Wait refresh time
		<-t
	}
}

// OnKeyPress registers a provided function to be called on a key press.
// Returns a function that when called unregisters the OnKeyPress function.
func OnKeyPress(key pixelgl.Button, f func(pixelgl.Button)) func() {
	onKeyPressedFuncsMutex.Lock()
	id := uuid.New().String()
	onKeyPressedFuncs[key] = map[string]func(pixelgl.Button){
		id: f,
	}
	onKeyPressedFuncsMutex.Unlock()

	return func() {
		onKeyPressedFuncsMutex.Lock()
		delete(onKeyPressedFuncs[key], id)
		onKeyPressedFuncsMutex.Unlock()
	}
}
