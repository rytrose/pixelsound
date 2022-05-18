package util

import "image/color"

// Uint8RGBA converts a color.Color to r, g, b, a represented as 0-255.
func Uint8RGBA(c color.Color) (r, g, b, a uint8) {
	ri, gi, bi, ai := c.RGBA()
	r = uint8(ri / 256)
	g = uint8(gi / 256)
	b = uint8(bi / 256)
	a = uint8(ai / 256)
	return
}
// FloatRGBA converts a color.Color to r, g, b, a represented as 0.0-1.0.
func FloatRGBA(c color.Color) (r, g, b, a float64) {
	ri, gi, bi, ai := c.RGBA()
	r = float64(ri) / 65535
	g = float64(gi) / 65535
	b = float64(bi) / 65535
	a = float64(ai) / 65535
	return
}
