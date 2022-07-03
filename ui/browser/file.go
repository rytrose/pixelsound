//go:build js

package browser

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/rytrose/pixelsound/log"
	"honnef.co/go/js/dom"
)

const audioInputID = "audioFile"
const imageInputID = "imageFile"

func readFilesFromInput(d dom.Document, id string) (chan []byte, error) {
	input := d.QuerySelector(fmt.Sprintf("#%s", id))
	if input == nil {
		return nil, fmt.Errorf("input with id %s does not exist", id)
	}
	inputEl, ok := input.(*dom.HTMLInputElement)
	if !ok {
		return nil, fmt.Errorf("element with id %s is not an input element", id)
	}

	fileChan := make(chan []byte)
	input.AddEventListener("change", false, func(e dom.Event) {
		go func() {
			fileObj := inputEl.Files()[0].Object
			fileObj.Call("arrayBuffer").
				Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					data := args[0]
					value := js.Global().Get("Uint8Array").New(data)
					b := make([]byte, value.Length())
					js.CopyBytesToGo(b, value)
					fileChan <- b
					return nil
				})).
				Call("catch", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
					err := args[0]
					log.Println("unable to get arrayBuffer from file input")
					println(err)
					return nil
				}))
		}()
	})
	return fileChan, nil
}

func (b *browser) readAudioFilesFromInput() {
	fileChan, err := readFilesFromInput(b.d, audioInputID)
	if err != nil {
		log.Fatalf("unable to read files from #%s input: %s", audioInputID, err)
		return
	}

	for file := range fileChan {
		b.setAudio(bytesReaderCloser{bytes.NewReader(file)}, "mp3")
		b.reloadCurrentMode()
	}
}
