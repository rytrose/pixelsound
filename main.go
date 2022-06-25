package main

import (
	"github.com/go-humble/detect"
	"github.com/rytrose/pixelsound/ui"
	"github.com/rytrose/pixelsound/ui/browser"
	"github.com/rytrose/pixelsound/ui/thick"
)

func main() {
	// Determine which UI to use
	var ui ui.UI
	if detect.IsBrowser() {
		ui = setupJS()
	} else {
		ui = setupDarwin()
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
