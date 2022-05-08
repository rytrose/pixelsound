//go:build js

package browser

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

// LoadImageFromURL loads a GIF/PNG/JPEG image given a URL to a valid file.
func LoadImageFromURL(url string) (image.Image, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	return image.Decode(resp.Body)
}
