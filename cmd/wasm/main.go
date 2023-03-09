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
	btnSpeed1 := getElementById("btn_speed_1")
	setOnClickHandler(btnSpeed1, func() { player.ChangeSpeed(0.5) })
	btnSpeed2 := getElementById("btn_speed_2")
	setOnClickHandler(btnSpeed2, func() { player.ChangeSpeed(1) })
	btnSpeed3 := getElementById("btn_speed_3")
	setOnClickHandler(btnSpeed3, func() { player.ChangeSpeed(1.5) })
	btnSpeed4 := getElementById("btn_speed_4")
	setOnClickHandler(btnSpeed4, func() { player.ChangeSpeed(2) })

	for {
		state := player.WaitForStateUpdate()
		drawCurrentState(state)
	}
	defer player.Close()

	<-stopChan
}
