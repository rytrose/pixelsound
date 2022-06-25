package player

import (
	"image"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/rytrose/pixelsound/api"
	"github.com/rytrose/pixelsound/util"
)

// Player controls audio playback of a PixelSound.
type Player struct {
	sr             beep.SampleRate   // Sample rate of playback
	bs             int               // Buffer size of playback
	i              image.Image       // Image being played
	ps             api.PixelSound    // Algorithms for traversal and sonification
	loc            image.Point       // Pixel location
	state          interface{}       // Previous state from sonification
	q              *Queue            // Streamer to queue up playback
	c              *beep.Ctrl        // Streamer to play/pause
	v              *effects.Volume   // Streamer to control volume
	PointChan      chan image.Point  // Writes the point being played
	usePointChan   bool              // If set, writes the point being played to PointChan
	LatestPoint    *image.Point      // The latest played point, access requires PointLock
	PointLock      util.PriorityLock // Lock for reading/writing the latest played point
	useLatestPoint bool              // If set, writes the point being played to LatestPoint
}

type PlayerOpt func(*Player)

func WithPointChan() PlayerOpt {
	return func(p *Player) {
		p.usePointChan = true
	}
}

func WithLatestPoint() PlayerOpt {
	return func(p *Player) {
		p.useLatestPoint = true
	}
}

// NewPlayer creates a Player.
func NewPlayer(sampleRate beep.SampleRate, bufferSize int, opts ...PlayerOpt) *Player {
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

	// Define Player
	p := &Player{
		sr: sampleRate,
		bs: bufferSize,
		q:  q,
		c:  c,
		v:  v,
		// Buffer so that very fast calls to PlayPixel don't get behind if the
		// reader is slow
		PointChan: make(chan image.Point, 60),
		PointLock: util.NewPriorityPreferenceLock(),
	}

	// Apply options
	for _, o := range opts {
		o(p)
	}

	return p
}

// SetImagePixelSound sets the current image and PixelSound.
func (p *Player) SetImagePixelSound(image image.Image, ps api.PixelSound) {
	p.i = image
	p.ps = ps
}

// Play plays a provided PixelSound for an image starting from provided coordinates.
func (p *Player) Play(image image.Image, ps api.PixelSound, start image.Point) {
	// Save playing image, PixelSound, and starting coordinates
	p.i = image
	p.ps = ps
	p.loc = start

	// Stop anything playing previously
	p.q.Clear()

	// Get the first pixel Streamer
	s, state := ps.Sonify(p.i.At(p.loc.X, p.loc.Y), p.sr, nil)
	p.state = state

	// Call the next pixel Streamer after the first is done
	n := beep.Seq(beep.Callback(func() {
		// Update with the first point
		p.updatePoint()
	}), s, beep.Callback(p.next))

	// Start playback by adding to mixer
	p.q.Add(n)
}

// next traverses the PixelSound and queues up the next pixel Streamer, if there is one.
func (p *Player) next() {
	var ok bool
	p.loc, ok = p.ps.Traverse(p.loc, p.i.Bounds())
	p.updatePoint()
	if ok {
		// Add this pixel Streamer, then the next
		s, state := p.ps.Sonify(p.i.At(p.loc.X, p.loc.Y), p.sr, p.state)
		p.state = state
		p.q.Add(beep.Seq(s, beep.Callback(p.next)))
	} else {
		// Add the final pixel Streamer
		s, state := p.ps.Sonify(p.i.At(p.loc.X, p.loc.Y), p.sr, p.state)
		p.state = state
		p.q.Add(s)
	}
}

// PlayPixel plays the pixel at the provided point.
func (p *Player) PlayPixel(point image.Point, queue bool) {
	p.loc = point
	p.updatePoint()
	s, state := p.ps.Sonify(p.i.At(point.X, point.Y), p.sr, p.state)
	p.state = state
	if !queue {
		p.q.Clear()
	}
	p.q.Add(s)
}

// updatePoint sends the currently playing point through PointChan,
// and/or updates LatestPoint.
func (p *Player) updatePoint() {
	if p.usePointChan {
		p.PointChan <- p.loc
	}
	if p.useLatestPoint {
		p.PointLock.Lock()
		p.LatestPoint = &p.loc
		p.PointLock.Unlock()
	}
}

// Stop clears the queue to stop playback.
func (p *Player) Stop() {
	p.q.Clear()
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
