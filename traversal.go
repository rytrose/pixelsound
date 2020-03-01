package main

import (
	"image"
	"math/rand"
)

// TtoBLtoR traverse an image top-to-bottom left-to-right
func TtoBLtoR(pX, pY int, b image.Rectangle) (int, int, bool) {

	if pX == b.Max.X-2 && pY == b.Max.Y-1 {
		return pX + 1, pY, false
	}
	if pX < b.Max.X-1 {
		return pX + 1, pY, true
	}
	return b.Min.X, pY + 1, true
}

// Random traverse an image in a random fashion.
func Random(pX, pY int, b image.Rectangle) (int, int, bool) {
	x := b.Min.X + rand.Intn(b.Max.X)
	y := b.Min.Y + rand.Intn(b.Max.Y)
	return x, y, true
}
