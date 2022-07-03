//go:build js

package browser

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/vincent-petithory/dataurl"
)

// loadAudioFromURL loads mp3/wav/ogg audio given a URL to a valid file.
// Requires the URL to have permissive CORS headers.
func loadAudioFromURL(url string) (io.ReadCloser, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", err
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	contentType := http.DetectContentType(d)
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
	return io.NopCloser(bytes.NewReader(d)), ext, nil
}

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

func (b *browser) setAudio(r bytesReaderCloser, ext string) {
	b.r = r
	b.ext = ext
}
