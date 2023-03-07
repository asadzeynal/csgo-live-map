package main

import (
	"log"

	"github.com/asadzeynal/csgo-live-map/engine"
)

func main() {
	stopChan := make(chan bool)
	fileInput := getElementById("file_input")

	file := <-setFileInputHandler(fileInput)

	player, err := engine.GetPlayer(file)
	if err != nil {
		log.Panic("error when getting player: %v", err)
	}

	drawMap(player.MapName)

	btnStop := getElementById("stop_button")
	setOnClickHandler(btnStop, player.Stop)
	btnPause := getElementById("pause_button")
	setOnClickHandler(btnPause, player.Pause)
	btnPlay := getElementById("play_button")
	setOnClickHandler(btnPlay, player.Play)
	btnNextRound := getElementById("next_round_button")
	setOnClickHandler(btnNextRound, player.NextRound)
	btnPrevRound := getElementById("prev_round_button")
	setOnClickHandler(btnPrevRound, player.PrevRound)

	for {
		state := player.WaitForStateUpdate()
		drawCurrentState(state)
	}
	defer player.Close()

	<-stopChan
}
