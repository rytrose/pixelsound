//go:build js

package browser

import (
	"image"

	"honnef.co/go/js/dom"
)

func OnMouseMove(el dom.Element, f func(image.Point)) func() {
	wrapper := el.AddEventListener("mousemove", false, func(e dom.Event) {
		x := e.Underlying().Get("offsetX").Int()
		y := e.Underlying().Get("offsetY").Int()
		go f(image.Point{x, y})
	})
	return func() {
		el.RemoveEventListener("mousemove", false, wrapper)
	}
}
