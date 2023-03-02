package engine

import (
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
)

type StateResult struct {
	Players []PlayerData
}

type PlayerData struct {
	Position Position
}

type Position struct {
	X float64
	Y float64
}

type engine struct {
	mapMetadata *Map
}

func (e *engine) getUsefulState(state demoinfocs.GameState) StateResult {
	players := state.Participants().Playing()
	playersData := make([]PlayerData, 0, len(players))

	for i := range players {
		p := players[i]
		vec := p.Position()
		// x, y := e.mapMetadata.TranslateScale(vec.X, vec.Y)
		x, y := vec.X, vec.Y

		pos := Position{X: x, Y: y}

		playersData = append(playersData, PlayerData{Position: pos})
	}

	return StateResult{Players: playersData}
}
