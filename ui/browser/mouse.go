//go:build js

package browser

import (
	"image"

	"honnef.co/go/js/dom"
)

func OnMouseMove(el dom.Element, f func(point image.Point, width int, height int)) func() {
	wrapper := el.AddEventListener("mousemove", false, func(e dom.Event) {
		width := e.CurrentTarget().Underlying().Get("width").Int()
		height := e.CurrentTarget().Underlying().Get("height").Int()
		x := e.Underlying().Get("offsetX").Int()
		y := e.Underlying().Get("offsetY").Int()
		go f(image.Point{x, y}, width, height)
	})
	return func() {
		el.RemoveEventListener("mousemove", false, wrapper)
	}
}
