package engine

import (
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
)

// Contains state for current tick
type StateResult struct {
	Players []PlayerData
}

// Contans current state for a single player
type PlayerData struct {
	Position      Position // Coordinates of player on 1024*1024 map image
	Team          byte     // 2 = T, 3 = CT
	IsAlive       bool     // true if player is alive
	ViewDirection float32  // 0-360 direction where player is looking on the 2d plane
}

// Position of an entity on the map
type Position struct {
	X float64
	Y float64
}

type engine struct {
	mapMetadata *Map
}

// Responsible for deriving useful state from demoinfocs.GameState and returning it
func (e *engine) getUsefulState(state demoinfocs.GameState) StateResult {
	players := state.Participants().Playing()
	playersData := make([]PlayerData, 0, len(players))

	for i := range players {
		pd := e.constructPlayerData(players[i])
		playersData = append(playersData, pd)
	}

	return StateResult{Players: playersData}
}

// Constructs and returns PlayerData object from demoinfocs.common.Player
// Uses TranslateScale function to translate from demo coordinates to 1024x1024 coordinates
func (e *engine) constructPlayerData(p *common.Player) PlayerData {
	demoPos := p.Position()
	posX, posY := e.mapMetadata.TranslateScale(demoPos.X, demoPos.Y)
	resPos := Position{X: posX, Y: posY}

	return PlayerData{
		Position:      resPos,
		Team:          byte(p.Team),
		IsAlive:       p.IsAlive(),
		ViewDirection: p.ViewDirectionX(),
	}
}
