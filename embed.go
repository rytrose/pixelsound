package main

import (
	"github.com/rytrose/pixelsound/log"

	"github.com/siongui/goef"
)

func generateEmbeddedAudio() {
	err := goef.GenerateGoPackage("browser", "ui/browser/audio/", "ui/browser/embedded_audio.go")
	if err != nil {
		log.Fatalf("unable to generate embedded audio files: %s", err)
	}
}
