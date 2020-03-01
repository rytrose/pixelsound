package main

import (
	"image"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

// Player controls audio playback of a PixelSound.
type Player struct {
	sr beep.SampleRate
	bs int
	i  image.Image
	ps PixelSound
	x  int
	y  int
	q  *Queue
	c  *beep.Ctrl
	v  *effects.Volume
}

// NewPlayer creates a Player.
func NewPlayer(sampleRate beep.SampleRate, bufferSize int) *Player {
	// Initialize the speaker
	speaker.Init(sampleRate, bufferSize)

	// Setup beep streamers
	q := &Queue{}
	c := &beep.Ctrl{
		Streamer: q,
		Paused:   false,
	}
	v := &effects.Volume{
		Streamer: c,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	// Start playing (plays silence until something is added)
	speaker.Play(q)

	// Return Player
	return &Player{
		sr: sampleRate,
		bs: bufferSize,
		q:  q,
		c:  c,
		v:  v,
	}
}

// Play plays a provided PixelSound for an image starting from provided coordinates.
func (p *Player) Play(image image.Image, ps PixelSound, x, y int) {
	// Save playing image, PixelSound, and starting coordinates
	p.i = image
	p.ps = ps
	p.x = x
	p.y = y

	// Stop anything playing previously
	p.q.Clear()

	// Get the first pixel Streamer
	s := ps.Sonify(p.i.At(p.x, p.y), p.sr)

	// Call the next pixel Streamer after the first is done
	n := beep.Seq(s, beep.Callback(p.next))

	// Start playback by adding to mixer
	p.q.Add(n)
}

// next traverses the PixelSound and queues up the next pixel Streamer, if there is one.
func (p *Player) next() {
	var ok bool
	p.x, p.y, ok = p.ps.Traverse(p.x, p.y, p.i.Bounds())
	pointChan <- image.Point{
		X: p.x,
		Y: p.y,
	}
	if ok {
		// Add this pixel Streamer, then the next
		p.q.Add(beep.Seq(p.ps.Sonify(p.i.At(p.x, p.y), p.sr), beep.Callback(p.next)))
	} else {
		// Add the final pixel Streamer
		p.q.Add(p.ps.Sonify(p.i.At(p.x, p.y), p.sr))
	}
}

// TogglePlayback toggles the playing/paused state of the player.
func (p *Player) TogglePlayback() {
	speaker.Lock()
	p.c.Paused = !p.c.Paused
	speaker.Unlock()
}

// Pause pauses the playback state of the player.
func (p *Player) Pause() {
	speaker.Lock()
	p.c.Paused = true
	speaker.Unlock()
}

// Resume resumes the playback state of the player.
func (p *Player) Resume() {
	speaker.Lock()
	p.c.Paused = false
	speaker.Unlock()
}

// Mute mutes the playback of the player.
func (p *Player) Mute() {
	speaker.Lock()
	p.v.Silent = true
	speaker.Unlock()
}

// Unmute unmutes the playback of the player.
func (p *Player) Unmute() {
	speaker.Lock()
	p.v.Silent = false
	speaker.Unlock()
}

// ToggleMute toggles the mute status of the playback of the player.
func (p *Player) ToggleMute() {
	speaker.Lock()
	p.v.Silent = !p.v.Silent
	speaker.Unlock()
}

// SetVolume sets the volume of the player.
// 0 is no volume change, negative numbers are quieter, positive numbers are louder.
func (p *Player) SetVolume(v float64) {
	speaker.Lock()
	p.v.Volume = v
	speaker.Unlock()
}
