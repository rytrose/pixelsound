package util

import (
	"time"

	"github.com/rytrose/pixelsound/log"
)

func TimeTrack(start time.Time, name string) time.Time {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
	return time.Now()
}
