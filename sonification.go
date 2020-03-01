package main

import (
	"image/color"
	"time"

	"github.com/faiface/beep"
)

// SineColor maps red to freq and green to duration.
func SineColor(c color.Color, sr beep.SampleRate) beep.Streamer {
	r, g, _, _ := c.RGBA()
	freq := 30 + (1600 * (float64(r) / 65535))
	durMS := 10 + (40 * (float64(g) / 65535))
	return beep.Take(sr.N(time.Duration(durMS)*time.Millisecond), Sine(sr, float64(freq)))
}
