package main

import (
	"log"
	"runtime"
	"strings"

	"github.com/rytrose/pixelsound/ui"
	"github.com/rytrose/pixelsound/ui/browser"
	"github.com/rytrose/pixelsound/ui/thick"
)

var isJS bool
var isDarwin bool

func main() {
	isJS = strings.Contains(runtime.GOOS, "js")
	isDarwin = strings.Contains(runtime.GOOS, "darwin")

	// Determine UI to use based on running os
	var ui ui.UI
	if isJS {
		ui = setupJS()
	} else if isDarwin {
		ui = setupDarwin()
	} else {
		log.Fatalf("Unsupported runtime OS: %s", runtime.GOOS)
	}

	// Run the UI
	ui.Run()
}

func setupJS() ui.UI {
	return browser.NewBrowser()
}

func setupDarwin() ui.UI {
	return thick.NewThickClient()
}
