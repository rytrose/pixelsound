//go:build !js

package thick

import (
	"sync"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"github.com/google/uuid"
)

type keyPressFunctionInstance struct {
	f                 func(pixelgl.Button)
	lastActivated     time.Time
	repeat            bool
	repeating         bool
	repeatDelay       time.Duration
	repeatSampleDelay time.Duration
}
type keyPressFunctionMap map[pixelgl.Button]map[string]*keyPressFunctionInstance

var keyMute sync.RWMutex
var onKeyPressedFuncMap = keyPressFunctionMap{}
var defaultRepeatDelay = 400 * time.Millisecond
var defaultRepeatSampleDelay = 100 * time.Millisecond

// KeyboardUpdate observes keyboard activity and calls registered functions.
func KeyboardUpdate(win *pixelgl.Window) {
	keyMute.Lock()
	// Check if registered keys were pressed
	for key := range onKeyPressedFuncMap {
		// Call function on key press
		if win.JustPressed(key) {
			for _, o := range onKeyPressedFuncMap[key] {
				o.lastActivated = time.Now()
				go o.f(key)
			}
		}
		// Handle repeat on press-and-hold
		if !win.JustPressed(key) && win.Pressed(key) {
			for _, o := range onKeyPressedFuncMap[key] {
				if !o.repeating {
					if time.Since(o.lastActivated) > o.repeatDelay {
						o.repeating = true
						o.lastActivated = time.Now()
						go o.f(key)
					}
				} else {
					if time.Since(o.lastActivated) > o.repeatSampleDelay {
						o.lastActivated = time.Now()
						go o.f(key)
					}
				}
			}
		}
		// Clear state on release
		if win.JustReleased(key) {
			for _, o := range onKeyPressedFuncMap[key] {
				o.repeating = false
			}
		}
	}
	keyMute.Unlock()
}

// OnKeyPress registers a provided function to be called on a key press.
// Returns a function that when called unregisters the OnKeyPress function.
func OnKeyPress(key pixelgl.Button, f func(pixelgl.Button), repeat bool) func() {
	id := uuid.New().String()
	keyMute.Lock()
	onKeyPressedFuncMap[key] = map[string]*keyPressFunctionInstance{
		id: {
			f:                 f,
			repeat:            repeat,
			repeating:         false,
			repeatDelay:       defaultRepeatDelay,
			repeatSampleDelay: defaultRepeatSampleDelay,
		},
	}
	keyMute.Unlock()

	return func() {
		keyMute.Lock()
		delete(onKeyPressedFuncMap[key], id)
		keyMute.Unlock()
	}
}
