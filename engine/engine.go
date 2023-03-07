package engine

import (
	"time"

	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
)

// Contains state for current tick
type StateResult struct {
	TeamT         Team
	TeamCt        Team
	RoundTimeLeft Second
	Nades         []Nade
	ScoreT        int
	ScoreCt       int
	CurrentRound  int
}

type Team struct {
	Players []PlayerData
	ClanTag string
}

type Second time.Duration

// Contans current state for a single player
type PlayerData struct {
	Name              string   // Player's in-game name
	Id                int      // Id 1...10 to be displayed on the map
	Position          Position // Coordinates of player on 1024*1024 map image
	LastAlivePosition Position // Position where player was last alive, used when isAlive == false
	Team              byte     // 2 = T, 3 = CT
	IsAlive           bool     // true if player is alive
	ViewDirection     float32  // 0-360 direction where player is looking on the 2d plane
	Kills             int
	Assists           int
	Deaths            int
	Money             int
	Equipped          string
	HP                int
}

type Nade struct {
	Positions []Position
	Type      string
}

// Position of an entity on the map
type Position struct {
	X float64
	Y float64
}

type engine struct {
	mapMetadata          *Map
	roundFreezeTimeEndAt time.Duration
	roundEndedAt         time.Duration
	playerIds            map[string]int // 1...10
	scoreT               int
	scoreCt              int
	currentRound         int
}

func (e *engine) constructPlayerIds(playersT []*common.Player, playersCt []*common.Player) {
	pIds := make(map[string]int)
	for i := range playersT {
		pIds[playersT[i].Name] = i + 1
	}
	for i := range playersCt {
		pIds[playersCt[i].Name] = i + 6
	}
	e.playerIds = pIds
}

// Responsible for deriving useful state from demoinfocs.GameState and returning it
func (e *engine) getUsefulState(state demoinfocs.GameState, currentTime time.Duration) StateResult {
	playersT := state.TeamTerrorists().Members()
	playersCt := state.TeamCounterTerrorists().Members()
	if e.playerIds == nil {
		e.constructPlayerIds(playersT, playersCt)
	}

	playersDataT := make([]PlayerData, 0, len(playersT))
	playersDataCt := make([]PlayerData, 0, len(playersCt))

	timeLimit, err := state.Rules().RoundTime()
	if err != nil {
		timeLimit = 115 * time.Nanosecond
	}

	roundTime := e.calculateRoundLeftTime(timeLimit, currentTime)

	for i := range playersT {
		pd := e.constructPlayerData(playersT[i])
		playersDataT = append(playersDataT, pd)
	}
	teamT := Team{
		Players: playersDataT,
		ClanTag: state.TeamTerrorists().ClanName(),
	}
	for i := range playersCt {
		pd := e.constructPlayerData(playersCt[i])
		playersDataCt = append(playersDataCt, pd)
	}
	teamCt := Team{
		Players: playersDataCt,
		ClanTag: state.TeamCounterTerrorists().ClanName(),
	}

	currentNades := state.GrenadeProjectiles()
	nadesRes := e.calculateNadeTrajectories(currentNades)

	return StateResult{
		TeamT:         teamT,
		TeamCt:        teamCt,
		RoundTimeLeft: roundTime,
		Nades:         nadesRes,
		ScoreT:        e.scoreT,
		ScoreCt:       e.scoreCt,
		CurrentRound:  e.currentRound,
	}
}

func (e *engine) calculateNadeTrajectories(nades map[int]*common.GrenadeProjectile) []Nade {
	res := make([]Nade, 0, len(nades))
	for i := range nades {
		nade := nades[i]
		trajectory := nade.Trajectory
		positions := make([]Position, 0, len(trajectory))
		for j := range trajectory {
			coords := trajectory[j]
			x, y := e.mapMetadata.TranslateScale(coords.X, coords.Y)
			positions = append(positions, Position{X: x, Y: y})
		}
		res = append(res, Nade{Positions: positions, Type: nade.WeaponInstance.String()})
	}
	return res
}

func (e *engine) calculateRoundLeftTime(timeLimit time.Duration, currentTime time.Duration) Second {
	if e.roundEndedAt >= e.roundFreezeTimeEndAt {
		return 0
	} else {
		return Second((timeLimit - (currentTime - e.roundFreezeTimeEndAt)) / time.Second)
	}
}

// Constructs and returns PlayerData object from demoinfocs.common.Player
// Uses TranslateScale function to translate from demo coordinates to 1024x1024 coordinates
func (e *engine) constructPlayerData(p *common.Player) PlayerData {
	demoPos := p.Position()
	posX, posY := e.mapMetadata.TranslateScale(demoPos.X, demoPos.Y)
	lastAlivePosX, lastAlivePosY := e.mapMetadata.TranslateScale(p.LastAlivePosition.X, p.LastAlivePosition.Y)

	var equipped string
	if p.ActiveWeapon() != nil {
		equipped = p.ActiveWeapon().String()
	}

	return PlayerData{
		Name:              p.Name,
		Id:                e.playerIds[p.Name],
		Position:          Position{X: posX, Y: posY},
		LastAlivePosition: Position{X: lastAlivePosX, Y: lastAlivePosY},
		Team:              byte(p.Team),
		IsAlive:           p.IsAlive(),
		ViewDirection:     p.ViewDirectionX(),
		Kills:             p.Kills(),
		Assists:           p.Assists(),
		Deaths:            p.Deaths(),
		Money:             p.Money(),
		Equipped:          equipped,
		HP:                p.Health(),
	}
}
