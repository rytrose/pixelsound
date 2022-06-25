//go:build js

package browser

import (
	"honnef.co/go/js/dom"
)

type keyCode int

const (
	keyLeft  = 37
	keyRight = 39
	keyUp    = 38
	keyDown  = 40
)

func OnKeyboardMove(w dom.Window, f func(keyCode)) func() {
	wrapper := w.AddEventListener("keydown", false, func(e dom.Event) {
		code := e.Underlying().Get("keyCode").Int()
		switch code {
		case keyLeft, keyRight, keyUp, keyDown:
			go f(keyCode(code))
		}
	})
	return func() {
		w.RemoveEventListener("keydown", false, wrapper)
	}
}
