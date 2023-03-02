package main

import (
	"bytes"
	"log"
	"syscall/js"
	"time"

	"github.com/asadzeynal/csgo-live-map/engine"
)

func main() {
	stopChan := make(chan bool)
	document := js.Global().Get("document")

	fileInput := document.Call("getElementById", "file_input")

	fileInput.Set("oninput", js.FuncOf(func(v js.Value, x []js.Value) any {
		fileInput.Get("files").Call("item", 0).Call("arrayBuffer").Call("then", js.FuncOf(func(v js.Value, x []js.Value) any {
			data := js.Global().Get("Uint8Array").New(x[0])
			dst := make([]byte, data.Get("length").Int())
			js.CopyBytesToGo(dst, data)

			player, err := engine.GetPlayer(bytes.NewReader(dst))
			if err != nil {
				log.Panic("error when getting player: %v", err)
			}

			ticker := time.NewTicker(500 * time.Millisecond)
			done := make(chan bool)

			go func() {
				for {
					select {
					case <-done:
						return
					case <-ticker.C:
						state := player.NextTick()
						js.Global().Call("drawPlayer", state)
					}
				}
			}()

			return nil
		}))

		return nil
	}))
	<-stopChan
}
