package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/asadzeynal/csgo-live-map/engine"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

func setOnClickHandler(element js.Value, f func()) {
	element.Set("onclick", js.FuncOf(func(v js.Value, x []js.Value) any {
		f()
		return nil
	}))
}

func setFileInputHandler(element js.Value) chan []byte {
	result := make(chan []byte)
	element.Set("oninput", js.FuncOf(func(v js.Value, x []js.Value) any {
		element.Get("files").
			Call("item", 0).
			Call("arrayBuffer").
			Call("then",
				js.FuncOf(func(v js.Value, x []js.Value) any {
					data := js.Global().Get("Uint8Array").New(x[0])
					dst := make([]byte, data.Get("length").Int())
					js.CopyBytesToGo(dst, data)
					result <- dst
					return nil
				}))
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
	js.Global().Call("updateState", string(stateJson))
}

func drawMap(mapName string) {
	js.Global().Call("drawMap", mapName)
}

func showError(msg string) {
	js.Global().Call("showError", msg)
}
