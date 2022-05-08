package api

import (
	"image"
	"image/color"

	"github.com/faiface/beep"
)

// TraverseFunc is a function that takes a 2D pixel location and 2D bounds and returns another 2D pixel location.
// If ok is true, there is more to traverse. If ok is false, the traversal is finished.
type TraverseFunc func(prev image.Point, bounds image.Rectangle) (next image.Point, ok bool)

// SonifyFunc is a function that takes a color and returns a beep.Streamer sonifying that color.
// It also takes an returns an interface for any state that should be passed between calls to
// the SonifyFunc.
type SonifyFunc func(color.Color, beep.SampleRate, interface{}) (beep.Streamer, interface{})

// PixelSound is an interface that describes how an image is traversed and sonified.
type PixelSound interface {
	Traverse(image.Point, image.Rectangle) (image.Point, bool)
	Sonify(color.Color, beep.SampleRate, interface{}) (beep.Streamer, interface{})
}

// PixelSounder is a struct that implements the PixelSound interface.
type PixelSounder struct {
	T TraverseFunc
	S SonifyFunc
}

// Traverse calls a TraverseFunc.
func (ps *PixelSounder) Traverse(prev image.Point, bounds image.Rectangle) (next image.Point, ok bool) {
	return ps.T(prev, bounds)
}

// Sonify calls a SonifyFunc.
func (ps *PixelSounder) Sonify(c color.Color, sr beep.SampleRate, state interface{}) (beep.Streamer, interface{}) {
	return ps.S(c, sr, state)
}
