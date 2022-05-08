package player

import (
	"github.com/faiface/beep"
)

// Queue plays streamers sequentially, and otherwise outputs silence.
type Queue struct {
	streamers []beep.Streamer
}

// Add adds a Streamer to the queue.
func (q *Queue) Add(streamers ...beep.Streamer) {
	q.streamers = append(q.streamers, streamers...)
}

// Clear removes all Streamers from the queue.
func (q *Queue) Clear() {
	q.streamers = q.streamers[:0]
}

// Stream streams the Streamer at the head of the queue, otherwise it streams silence.
func (q *Queue) Stream(samples [][2]float64) (n int, ok bool) {
	// We use the filled variable to track how many samples we've
	// successfully filled already. We loop until all samples are filled.
	filled := 0
	for filled < len(samples) {
		// There are no streamers in the queue, so we stream silence.
		if len(q.streamers) == 0 {
			for i := range samples[filled:] {
				samples[i][0] = 0
				samples[i][1] = 0
			}
			break
		}

		// We stream from the first streamer in the queue.
		n, ok := q.streamers[0].Stream(samples[filled:])
		// If it's drained, we pop it from the queue, thus continuing with
		// the next streamer.
		if !ok {
			q.streamers = q.streamers[1:]
		}
		// We update the number of filled samples.
		filled += n
	}
	return len(samples), true
}

// Err returns no error.
func (q *Queue) Err() error {
	return nil
}
