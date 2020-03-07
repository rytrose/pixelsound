package main

import (
	"math"
)

// Sine is a simple Sine wave Streamer.
type Sine struct {
	Freq  float64
	Amp   float64
	phase float64
	SR    float64
}

// NewSine is a Sine factory
func NewSine(freq float64, amp float64, sr float64) *Sine {
	return &Sine{
		Freq: freq,
		Amp:  amp,
		SR:   sr,
	}
}

// Stream returns samples of the sine wave.
func (s *Sine) Stream(samples [][2]float64) (n int, ok bool) {
	for i := range samples {
		if s.phase >= 1.0 {
			s.phase -= 1.0
		}
		y := math.Sin(2*math.Pi*s.phase) * s.Amp
		samples[i][0] = y
		samples[i][1] = y
		s.phase += (s.Freq / s.SR)
	}
	return len(samples), true
}

// Err returns no error.
func (s *Sine) Err() error {
	return nil
}
