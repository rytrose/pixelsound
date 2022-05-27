package traversal

import (
	"image"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

// Random traverses an image in a random fashion.
func Random(prev image.Point, bounds image.Rectangle) (image.Point, bool) {
	x := bounds.Min.X + rand.Intn(bounds.Max.X)
	y := bounds.Min.Y + rand.Intn(bounds.Max.Y)
	return image.Point{x, y}, true
}
