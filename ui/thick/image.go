package thick

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
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
func DrawImage(win *pixelgl.Window, sprite *pixel.Sprite, imd *imdraw.IMDraw, color color.Color, point image.Point) {
	// image has (0, 0) be top left, pixel has (0, 0) be bottom left
	pt := convertImageToPixel(win, point)

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

// convertImageToPixel converts coordinates from image.Point to pixel.Vec.
// pixel considers Y=0 to be bottom-left, while image considers Y=0 to
// be top-left.
func convertImageToPixel(win *pixelgl.Window, point image.Point) pixel.Vec {
	h := win.Bounds().H()
	newY := h - float64(point.Y)
	return pixel.V(float64(point.X), newY)
}

// convertPixelToImage converts coordinates from pixel.Vec to image.Point.
// pixel considers Y=0 to be bottom-left, while image considers Y=0 to
// be top-left.
func convertPixelToImage(win *pixelgl.Window, p pixel.Vec) image.Point {
	h := win.Bounds().H()
	newY := math.Abs(p.Y - h)
	return image.Point{int(p.X), int(newY)}
}
