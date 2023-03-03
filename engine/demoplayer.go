package engine

import (
	"fmt"
	"io"
	"time"

	"github.com/markus-wa/demoinfocs-golang/msg"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
)

type DemoPlayer struct {
	mapName       string
	parser        demoinfocs.Parser
	e             *engine
	playbackSpeed float64
	tickRate      float64
	ticker        *time.Ticker
}

var player *DemoPlayer = nil

func (dp *DemoPlayer) Close() {
	if dp != nil && dp.parser != nil {
		dp.parser.Close()
	}
	player = nil
}

func (dp *DemoPlayer) Pause() {
	dp.ticker.Stop()
}

func (dp *DemoPlayer) Play() {
	dp.refreshTicker()
}

func (dp *DemoPlayer) refreshTicker() {
	dp.ticker.Reset(time.Second / time.Duration(dp.tickRate) * time.Duration(dp.playbackSpeed))
}

func (dp *DemoPlayer) ChangeSpeed(speed float64) {
	dp.playbackSpeed = speed
	dp.refreshTicker()
}

func (dp *DemoPlayer) NextTick() StateResult {
	if dp.parser.Progress() == 1 {
		return StateResult{}
	}
	dp.parser.ParseNextFrame()
	res := dp.e.getUsefulState(dp.parser.GameState())

	return res
}

func GetPlayer(file io.Reader) (*DemoPlayer, error) {
	if player != nil {
		return player, nil
	}

	p := demoinfocs.NewParser(file)
	header, err := p.ParseHeader()
	if err != nil {
		return nil, fmt.Errorf("unable to parse demo headers: %v", err)
	}

	mapName := header.MapName

	if mapName != "de_ancient" {
		return nil, fmt.Errorf("only de_ancient is supported now")
	}

	var mapMetadata Map
	p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_ServerInfo) {
		fmt.Println(mapName, msg.MapCrc)
		// Get metadata for the map that the game was played on for coordinate translations
		mapMetadata = GetMapMetadata(mapName, msg.GetMapCrc())
		fmt.Println(mapMetadata)
	})

	for !p.GameState().IsMatchStarted() {
		p.ParseNextFrame()
	}

	e := engine{
		mapMetadata: &mapMetadata,
	}

	tickRate := p.TickRate()
	speed := 1.0

	ticker := time.NewTicker(time.Second / time.Duration(tickRate) * time.Duration(speed))

	player = &DemoPlayer{
		mapName:       mapName,
		parser:        p,
		e:             &e,
		playbackSpeed: speed,
		tickRate:      tickRate,
		ticker:        ticker,
	}

	return player, nil
}
