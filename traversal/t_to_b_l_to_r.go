package traversal

import (
	"image"
)

// TtoBLtoR traverses an image top-to-bottom left-to-right.
func TtoBLtoR(prev image.Point, bounds image.Rectangle) (image.Point, bool) {
	if prev.X == bounds.Max.X-2 && prev.Y == bounds.Max.Y-1 {
		return image.Point{prev.X + 1, prev.Y}, false
	}
	if prev.X < bounds.Max.X-1 {
		return image.Point{prev.X + 1, prev.Y}, true
	}
	return image.Point{bounds.Min.X, prev.Y + 1}, true
}
