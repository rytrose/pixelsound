package main

import (
	"math"

	"github.com/faiface/beep"
)

// Sine is a simple Sine wave Streamer.
func Sine(sr beep.SampleRate, freq float64, t float64) beep.Streamer {
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			y := math.Sin(math.Pi * freq * t)
			samples[i][0] = y
			samples[i][1] = y
			t += sr.D(1).Seconds()
		}
		return len(samples), true
	})
}
