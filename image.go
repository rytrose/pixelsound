package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/faiface/pixel/imdraw"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// LoadImageFromFile loads a GIF/PNG/JPEG image given a path to a file.
func LoadImageFromFile(path string) (image.Image, string, error) {
	reader, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer reader.Close()

	return image.Decode(reader)
}

// DrawImage draws an image with the currently "playing" pixel enlarged.
func DrawImage(win *pixelgl.Window, sprite *pixel.Sprite, imd *imdraw.IMDraw, color color.Color, x int, y int) {
	// image has (0, 0) be top left, pixel has (0, 0) be bottom left
	pt := convertImageToPixel(win, x, y)

	// Clear IMDraw
	imd.Clear()

	// Draw picture
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	// Update imd with colored pixel
	imd.Color = color
	imd.Push(pixel.V(pt.X-5, pt.Y-5), pixel.V(pt.X+5, pt.Y+5))
	imd.Rectangle(0)

	// Trace around rectangle
	imd.Color = pixel.RGB(0, 0, 0)
	imd.Push(pixel.V(pt.X-5, pt.Y-5), pixel.V(pt.X+5, pt.Y+5))
	imd.Rectangle(1)

	// Draw rectangle
	imd.Draw(win)
}

// convertImageToPixel converts coordinates from image to pixel
func convertImageToPixel(win *pixelgl.Window, x int, y int) pixel.Vec {
	h := win.Bounds().H()
	newY := h - float64(y)
	return pixel.V(float64(x), newY)
}

// convertPixelToImage converts coordinates from pixel to image
func convertPixelToImage(win *pixelgl.Window, p pixel.Vec) pixel.Vec {
	h := win.Bounds().H()
	newY := math.Abs(p.Y - h)
	return pixel.V(p.X, newY)
}

// hexColor returns an HTML hex-representation of c. The alpha channel is dropped
// and precision is truncated to 8 bits per channel
func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
