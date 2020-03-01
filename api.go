package main

import (
	"image"
	"image/color"

	"github.com/faiface/beep"
)

// TraverseFunc is a function that takes a 2D pixel location and 2D bounds and returns another 2D pixel location.
// If ok is true, there is more to traverse. If ok is false, the traversal is finished.
type TraverseFunc func(prevX, prevY int, bounds image.Rectangle) (x, y int, ok bool)

// SonifyFunc is a function that takes a color and returns a beep.Streamer sonifying that color.
type SonifyFunc func(color.Color, beep.SampleRate) beep.Streamer

// PixelSound is an interface that describes how an image is traversed and sonified.
type PixelSound interface {
	Traverse(prevX, prevY int, b image.Rectangle) (x, y int, ok bool)
	Sonify(c color.Color, sr beep.SampleRate) (s beep.Streamer)
}

// PixelSounder is a struct that implements the PixelSound interface.
type PixelSounder struct {
	T TraverseFunc
	S SonifyFunc
}

// Traverse calls a TraverseFunc.
func (ps *PixelSounder) Traverse(prevX, prevY int, b image.Rectangle) (x, y int, ok bool) {
	return ps.T(prevX, prevY, b)
}

// Sonify calls a SonifyFunc.
func (ps *PixelSounder) Sonify(c color.Color, sr beep.SampleRate) beep.Streamer {
	return ps.S(c, sr)
}
