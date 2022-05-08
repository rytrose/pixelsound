//go:build js

package thick

import "github.com/rytrose/pixelsound/ui"

// Returns a nil thick client UI when compiling for a non-darwin OS.
func NewThickClient() ui.UI {
	return nil
}
