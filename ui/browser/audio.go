//go:build js

package browser

import (
	"fmt"
	"net/http"

	"github.com/vincent-petithory/dataurl"
)

func decodeAudioFromDataURL(s string) ([]byte, string, error) {
	dataURL, err := dataurl.DecodeString(s)
	if err != nil {
		return nil, "", err
	}
	contentType := http.DetectContentType(dataURL.Data)
	ext := ""
	switch contentType {
	case "audio/mpeg":
		ext = "mp3"
	case "audio/wave":
		ext = "wav"
	case "application/ogg":
		ext = "ogg"
	}
	if ext == "" {
		return nil, "", fmt.Errorf("unable to handle content-type: %s", contentType)
	}
	return dataURL.Data, ext, nil
}
