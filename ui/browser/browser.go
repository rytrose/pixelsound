//go:build !js

package browser

import "github.com/rytrose/pixelsound/ui"

// Returns a nil browser UI when compiling for a non-JS OS.
func NewBrowser() ui.UI {
	return nil
}
