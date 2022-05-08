package sonification

import (
	"image/color"
	"io"
	"log"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
	"github.com/rytrose/pixelsound/api"
	"github.com/rytrose/pixelsound/util"
)

// resampleQuality determines the quality when resampling.
const resampleQuality = 3

// NewAudioScrubber returns a SonifyFunc that uses RGB to determine playback location, number of samples, and speed
// of the provided audio buffer formatted with the provided extension.
func NewAudioScrubber(r io.ReadCloser, ext string) api.SonifyFunc {
	var audioStreamer beep.StreamSeekCloser
	var audioFormat beep.Format
	var err error
	switch ext {
	case "mp3":
		audioStreamer, audioFormat, err = mp3.Decode(r)
	case "wav":
		audioStreamer, audioFormat, err = wav.Decode(r)
	case "ogg":
		audioStreamer, audioFormat, err = vorbis.Decode(r)
	case "flac":
		audioStreamer, audioFormat, err = flac.Decode(r)
	default:
		log.Fatalf("unable to decode audio file with extension %s", ext)
	}
	if err != nil {
		log.Fatalf("unable to decode audio: %s", err)
	}
	audioBuffer := beep.NewBuffer(audioFormat)
	audioBuffer.Append(audioStreamer)
	audioStreamer.Close()

	return func(c color.Color, sr beep.SampleRate, state interface{}) (beep.Streamer, interface{}) {
		return audioScrubber(audioBuffer, audioStreamer, c, sr, state)
	}
}

// audioScrubber uses RGB to determine audio buffer playback location, number of samples, and speed.
func audioScrubber(audioBuffer *beep.Buffer, audioStreamer beep.StreamSeekCloser, c color.Color, sr beep.SampleRate, state interface{}) (beep.Streamer, interface{}) {
	r, g, b, _ := util.FloatRGBA(c)
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
