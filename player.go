package main

import (
	"image"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// Player controls audio playback of a PixelSound.
type Player struct {
	win *pixelgl.Window // Window of GUI
	sr  beep.SampleRate // Sample rate of playback
	bs  int             // Buffer size of playback
	i   image.Image     // Image being played
	ps  PixelSound      // Algorithms for traversal and sonification
	x   int             // Previous pixel x
	y   int             // Previous pixel y
	st  interface{}     // Previous state from sonify
	q   *Queue          // Streamer to queue up playback
	c   *beep.Ctrl      // Streamer to play/pause
	v   *effects.Volume // Streamer to control volume
}

// NewPlayer creates a Player.
func NewPlayer(win *pixelgl.Window, sampleRate beep.SampleRate, bufferSize int) *Player {
	// Initialize the speaker
	speaker.Init(sampleRate, bufferSize)

	// Initialize sonification algorithms
	InitSonification(sampleRate)

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
		win: win,
		sr:  sampleRate,
		bs:  bufferSize,
		q:   q,
		c:   c,
		v:   v,
	}
}

// SetImagePixelSound sets the current image and PixelSound.
func (p *Player) SetImagePixelSound(image image.Image, ps PixelSound) {
	p.i = image
	p.ps = ps
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
	s, st := ps.Sonify(p.i.At(p.x, p.y), p.sr, nil)
	p.st = st

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
		s, st := p.ps.Sonify(p.i.At(p.x, p.y), p.sr, p.st)
		p.st = st
		p.q.Add(beep.Seq(s, beep.Callback(p.next)))
	} else {
		// Add the final pixel Streamer
		s, st := p.ps.Sonify(p.i.At(p.x, p.y), p.sr, p.st)
		p.st = st
		p.q.Add(s)
	}
}

// PlayPixel plays the pixel at the provided point.
func (p *Player) PlayPixel(po pixel.Vec) {
	po = po.Floor()
	po = convertPixelToImage(p.win, po)
	pointChan <- image.Point{
		X: int(po.X),
		Y: int(po.Y),
	}
	s, st := p.ps.Sonify(p.i.At(int(po.X), int(po.Y)), p.sr, p.st)
	p.st = st
	p.q.Add(s)
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
