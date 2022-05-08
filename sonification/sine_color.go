package sonification

import (
	"image/color"
	"time"

	"github.com/faiface/beep"
	"github.com/rytrose/pixelsound/api"
	"github.com/rytrose/pixelsound/util"
)

// NewSineColor returns a SonifyFunc that maps red to freq and green to duration of
// a sine wave.
func NewSineColor(sr beep.SampleRate) api.SonifyFunc {
	sine := NewSine(440, 1.0, float64(sr.N(1*time.Second)))
	return func(c color.Color, sr beep.SampleRate, state interface{}) (beep.Streamer, interface{}) {
		return sineColor(sine, c, sr, state)
	}
}

// sineColor maps red to freq and green to duration.
func sineColor(sine *Sine, c color.Color, sr beep.SampleRate, state interface{}) (beep.Streamer, interface{}) {
	r, g, _, _ := util.FloatRGBA(c)
	sine.Freq = 30 + (1600 * r)
	durMS := 10 + (40 * g)

	return beep.Take(sr.N(time.Duration(durMS)*time.Millisecond), sine), nil
}
