package engine

import (
	"time"

	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"
)

type engine struct {
	mapMetadata          *Map
	roundFreezeTimeEndAt time.Duration
	roundEndedAt         time.Duration
	playerIds            map[string]int // 1...10
	scoreT               int
	scoreCt              int
	currentRound         int
	isBombPlanted        bool
	isBombDefused        bool
	bombPlantedAt        time.Duration
}

// Constructs local player ids [1...10], called once.
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
		Score:   e.scoreT,
		ClanTag: state.TeamTerrorists().ClanName(),
	}
	for i := range playersCt {
		pd := e.constructPlayerData(playersCt[i])
		playersDataCt = append(playersDataCt, pd)
	}
	teamCt := Team{
		Players: playersDataCt,
		Score:   e.scoreCt,
		ClanTag: state.TeamCounterTerrorists().ClanName(),
	}

	nades := e.calculateNadeTrajectories(state.GrenadeProjectiles())
	infernos := e.calculateInfernosBorders(state.Infernos())

	bomb := e.getBomb(state.Bomb())

	return StateResult{
		TeamT:         teamT,
		TeamCt:        teamCt,
		RoundTimeLeft: roundTime,
		Nades:         nades,
		Infernos:      infernos,
		CurrentRound:  e.currentRound,
		Bomb:          bomb,
	}
}

func (e *engine) getBomb(bomb *common.Bomb) Bomb {
	carrier := bomb.Carrier
	bombPos := bomb.Position()
	var carrierId int
	if carrier != nil {
		carrierId = e.playerIds[carrier.Name]
	}
	x, y := e.mapMetadata.TranslateScale(bombPos.X, bombPos.Y)

	return Bomb{
		CarrierId: carrierId,
		Position:  Position{X: x, Y: y},
		isPlanted: e.isBombPlanted,
		isDefused: e.isBombDefused,
	}
}

func (e *engine) calculateInfernosBorders(infernos map[int]*common.Inferno) []Inferno {
	res := make([]Inferno, 0, len(infernos))
	for k := range infernos {
		inferno := infernos[k]
		fires := inferno.Fires().Active().ConvexHull2D()
		borders := make([]Position, 0, len(fires))
		for i := range fires {
			coords := fires[i]
			x, y := e.mapMetadata.TranslateScale(coords.X, coords.Y)
			borders = append(borders, Position{X: x, Y: y})
		}
		res = append(res, Inferno{BorderPositions: borders})
	}
	return res
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
	if e.isBombPlanted {
		return Second(((time.Second * 45) - (currentTime - e.bombPlantedAt)) / time.Second)
	}
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

	flashedUntil := p.FlashDurationTimeRemaining() / time.Millisecond

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
		FlashTimeLeft:     Millisecond(flashedUntil),
	}
}
