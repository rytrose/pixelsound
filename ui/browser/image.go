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

// loadImageFromURL loads a GIF/PNG/JPEG image given a URL to a valid file.
// Requires the URL to have permissive CORS headers.
func loadImageFromURL(url string) (image.Image, *dataurl.DataURL, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return decodeImage(b)
}

func decodeImage(b []byte) (image.Image, *dataurl.DataURL, error) {
	im, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return nil, nil, err
	}
	dataURL := dataurl.New(b, http.DetectContentType(b))
	return im, dataURL, nil
}

func decodeImageFromDataURL(s string) (image.Image, error) {
	dataURL, err := dataurl.DecodeString(s)
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(bytes.NewReader(dataURL.Data))
	if err != nil {
		return nil, err
	}
	return im, nil
}
