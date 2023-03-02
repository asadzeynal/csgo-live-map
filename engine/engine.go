package engine

import "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"

type IncomingState struct {
	Players []*common.Player
}

type StateResult struct {
	Players []PlayerData
}

type PlayerData struct {
	Position Position
}

type Position struct {
	X float64
	Y float64
	Z float64
}

func processOneTick(data IncomingState) StateResult {
	players := data.Players
	playersData := make([]PlayerData, 0, len(players))

	for i := range players {
		p := players[i]
		vec := p.Position()
		pos := Position{X: vec.X, Y: vec.Y, Z: vec.Z}

		playersData = append(playersData, PlayerData{Position: pos})
	}

	return StateResult{Players: playersData}
}
