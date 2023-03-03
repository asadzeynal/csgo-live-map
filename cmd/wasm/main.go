package main

import (
	"log"

	"github.com/asadzeynal/csgo-live-map/engine"
)

func main() {
	stopChan := make(chan bool)
	fileInput := getElementById("file_input")
	playButton := getElementById("play_button")
	setOnclickHandler(playButton)

	fileReaderChan := setFileInputHandler(fileInput)
	fileReader := <-fileReaderChan

	player, err := engine.GetPlayer(fileReader)
	if err != nil {
		log.Panic("error when getting player: %v", err)
	}

	for {
		state := player.WaitForStateUpdate()
		drawCurrentState(state)
	}
	defer player.Close()

	<-stopChan
}
