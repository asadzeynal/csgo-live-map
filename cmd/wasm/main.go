package main

import (
	"log"

	"github.com/asadzeynal/csgo-live-map/engine"
)

func main() {
	stopChan := make(chan bool)
	fileInput := getElementById("file_input")

	fileReaderChan := setFileInputHandler(fileInput, "oninput")
	fileReader := <-fileReaderChan

	player, err := engine.GetPlayer(fileReader)
	if err != nil {
		log.Panic("error when getting player: %v", err)
	}
	defer player.Close()

	<-stopChan
}
