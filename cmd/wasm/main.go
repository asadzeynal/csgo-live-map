package main

import (
	"bytes"
	"log"

	"github.com/asadzeynal/csgo-live-map/engine"
)

func main() {
	stopChan := make(chan bool)
	fileInput := getElementById("file_input")
	playButton := getElementById("play_button")

	file := <-setFileInputHandler(fileInput)
	player, err := engine.GetPlayer(bytes.NewReader(file))
	if err != nil {
		log.Panic("error when getting player: %v", err)
	}

	drawMap(player.MapName)

	setOnclickHandler(playButton, player.PlayPause)

	for {
		state := player.WaitForStateUpdate()
		drawCurrentState(state)
	}
	defer player.Close()

	<-stopChan
}
