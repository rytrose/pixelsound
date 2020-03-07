package main

import (
	"image/color"
	"time"

	"github.com/faiface/beep"
)

var sine *Sine

// InitSonification populates Streamers.
func InitSonification(sr beep.SampleRate) {
	sine = NewSine(440, 1.0, float64(sr.N(1*time.Second)))
}

// SineColor maps red to freq and green to duration.
func SineColor(c color.Color, sr beep.SampleRate, st interface{}) (beep.Streamer, interface{}) {
	r, g, _, _ := FloatRGBA(c)
	sine.Freq = 30 + (1600 * r)
	durMS := 10 + (40 * g)

	return beep.Take(sr.N(time.Duration(durMS)*time.Millisecond), sine), nil
}
