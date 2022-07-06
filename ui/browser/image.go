//go:build js

package browser

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/vincent-petithory/dataurl"
)

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
