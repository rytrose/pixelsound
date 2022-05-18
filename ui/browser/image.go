//go:build js

package browser

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"

	"github.com/vincent-petithory/dataurl"
)

// LoadImageFromURL loads a GIF/PNG/JPEG image given a URL to a valid file.
func LoadImageFromURL(url string) (image.Image, *dataurl.DataURL, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	im, _, err := image.Decode(bytes.NewReader(d))
	if err != nil {
		return nil, nil, err
	}
	dataURL := dataurl.New(d, http.DetectContentType(d))
	return im, dataURL, nil
}
