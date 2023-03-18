package main

import (
	"errors"
	"fmt"

	"github.com/asadzeynal/csgo-live-map/engine"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
)

func main() {
	stopChan := make(chan bool)

	player := waitForPlayer()
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

func waitForPlayer() *engine.DemoPlayer {
	err := errors.New("")
	var player *engine.DemoPlayer
	for err != nil {
		player, err = initPlayer()
		if err != nil {
			fmt.Println(err)
			if errors.As(err, &demoinfocs.ErrInvalidFileType) {
				showError("Invalid file type")
				continue
			}
			showError(err.Error())
		}
	}
	return player
}

func initPlayer() (*engine.DemoPlayer, error) {
	fileInput := getElementById("file_input")
	file := <-setFileInputHandler(fileInput)

	player, err := engine.GetPlayer(file)
	if err != nil {
		return nil, err
	}

	return player, nil
}
