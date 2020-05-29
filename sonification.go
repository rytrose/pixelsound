package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

const resampleQuality = 3

var (
	sine          *Sine
	audioFilename string
	audioStreamer beep.StreamSeekCloser
	audioFormat   beep.Format
	audioBuffer   *beep.Buffer
)

var sonifyFuncs = map[string]SonifyFunc{
	"SineColor":     SineColor,
	"AudioScrubber": AudioScrubber,
}

// InitSonification populates Streamers.
func InitSonification(sr beep.SampleRate) {
	sine = NewSine(440, 1.0, float64(sr.N(1*time.Second)))
	if audioFilename != "" {
		f, err := os.Open(audioFilename)
		if err != nil {
			panic(fmt.Sprintf("unable to open file %s: %s", audioFilename, err))
		}
		split := strings.Split(audioFilename, ".")
		ext := strings.ToLower(split[len(split)-1])
		switch ext {
		case "mp3":
			audioStreamer, audioFormat, err = mp3.Decode(f)
		case "wav":
			audioStreamer, audioFormat, err = wav.Decode(f)
		case "ogg":
			audioStreamer, audioFormat, err = vorbis.Decode(f)
		case "flac":
			audioStreamer, audioFormat, err = flac.Decode(f)
		default:
			panic(fmt.Sprintf("unable to decode audio file with extension %s", ext))
		}
		if err != nil {
			panic(fmt.Sprintf("unable to open %s %s: %s", ext, audioFilename, err))
		}
		audioBuffer = beep.NewBuffer(audioFormat)
		audioBuffer.Append(audioStreamer)
		audioStreamer.Close()
	}
}

// SineColor maps red to freq and green to duration.
func SineColor(c color.Color, sr beep.SampleRate, st interface{}) (beep.Streamer, interface{}) {
	r, g, _, _ := FloatRGBA(c)
	sine.Freq = 30 + (1600 * r)
	durMS := 10 + (40 * g)

	return beep.Take(sr.N(time.Duration(durMS)*time.Millisecond), sine), nil
}

// AudioScrubber uses RGB to determine audio buffer playback location, number of samples, and speed.
func AudioScrubber(c color.Color, sr beep.SampleRate, st interface{}) (beep.Streamer, interface{}) {
	r, g, b, _ := FloatRGBA(c)
	bufferLen := audioBuffer.Len()
	// G - duration
	maxDurPercentage := 0.1
	durSamples := int(maxDurPercentage * g * float64(bufferLen))
	// R - location
	startSample := int(r * float64(bufferLen-durSamples))
	s := audioBuffer.Streamer(startSample, startSample+durSamples)
	// B - speed (resampling ratio)
	var ratio float64
	if b < 0.5 {
		// Speed up when less blue
		ratio = 1 + (2.0 * b)
	} else {
		// Slow down when more blue
		ratio = 0.5 + (1 - b)
	}
	resampledStreamer := beep.ResampleRatio(resampleQuality, ratio, s)

	return resampledStreamer, nil
}
