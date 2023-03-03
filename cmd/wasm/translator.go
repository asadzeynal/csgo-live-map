package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/asadzeynal/csgo-live-map/engine"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

func setOnclickHandler(element js.Value) {
	element.Set("onclick", js.FuncOf(func(v js.Value, x []js.Value) any {
		fmt.Println("button clicked")
		return nil
	}))
}

func setFileInputHandler(element js.Value) chan *bytes.Reader {
	result := make(chan *bytes.Reader)
	element.Set("oninput", js.FuncOf(func(v js.Value, x []js.Value) any {
		element.Get("files").Call("item", 0).Call("arrayBuffer").Call("then", func(v js.Value, x []js.Value) any {
			data := js.Global().Get("Uint8Array").New(x[0])
			dst := make([]byte, data.Get("length").Int())
			js.CopyBytesToGo(dst, data)
			result <- bytes.NewReader(dst)
			return nil
		})
		return nil
	}))
	return result
}

func getElementById(id string) js.Value {
	return document.Call("getElementById", id)
}

func drawCurrentState(state engine.StateResult) {
	stateJson, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}

	js.Global().Call("drawPlayer", string(stateJson))
}

// func f(v js.Value, x []js.Value) any {
// 	ticker := time.NewTicker(500 * time.Millisecond)
// 	done := make(chan bool)

// 	go func() {
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case <-ticker.C:
// 				state := player.NextTick()
// 				js.Global().Call("drawPlayer", state)
// 			}
// 		}
// 	}()

// 	return nil

// }
