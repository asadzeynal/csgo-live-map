package engine

import (
	"fmt"
	"io"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/msg"
)

var supportedMaps map[string]struct{} = map[string]struct{}{
	"de_ancient":  {},
	"de_cache":    {},
	"de_dust2":    {},
	"de_inferno":  {},
	"de_mirage":   {},
	"de_nuke":     {},
	"de_overpass": {},
	"de_train":    {},
	"de_vertigo":  {},
}

type DemoPlayer struct {
	IsPaused      bool
	mapName       string
	parser        demoinfocs.Parser
	e             *engine
	playbackSpeed float64
	tickRate      float64
	ticker        *time.Ticker
	result        chan StateResult
}

var player *DemoPlayer = nil

func (dp *DemoPlayer) Close() {
	if dp != nil && dp.parser != nil {
		dp.parser.Close()
	}
	player = nil
}

func (dp *DemoPlayer) PlayPause() {
	if dp.IsPaused {
		dp.Play()
		return
	}
	dp.Pause()
}

func (dp *DemoPlayer) Pause() {
	if dp.IsPaused {
		return
	}
	dp.ticker.Stop()
	dp.IsPaused = true
}

func (dp *DemoPlayer) Play() {
	if !dp.IsPaused {
		return
	}
	dp.refreshTicker()
	dp.IsPaused = false
}

func (dp *DemoPlayer) refreshTicker() {
	dp.ticker.Reset(time.Second / time.Duration(dp.tickRate) * time.Duration(dp.playbackSpeed))
}

func (dp *DemoPlayer) ChangeSpeed(speed float64) {
	dp.playbackSpeed = speed
	dp.refreshTicker()
}

func (dp *DemoPlayer) nextTick() StateResult {
	if dp.parser.Progress() == 1 {
		return StateResult{}
	}
	dp.parser.ParseNextFrame()
	res := dp.e.getUsefulState(dp.parser.GameState())

	return res
}

func (dp *DemoPlayer) WaitForStateUpdate() StateResult {
	return <-dp.result
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

	if !isMapSupported(mapName) {
		return nil, fmt.Errorf("map %v is not supported right now", mapName)
	}

	player = &DemoPlayer{
		mapName:       mapName,
		parser:        p,
		IsPaused:      true,
		playbackSpeed: 1.0,
		result:        make(chan StateResult),
	}

	e := engine{}
	p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_ServerInfo) {
		mmd := GetMapMetadata(mapName, msg.GetMapCrc())
		e.mapMetadata = &mmd
	})

	for !p.GameState().IsMatchStarted() {
		p.ParseNextFrame()
	}
	tickRate := p.TickRate()

	ticker := time.NewTicker(time.Second / time.Duration(tickRate) * time.Duration(player.playbackSpeed))
	ticker.Stop()

	player.e = &e
	player.tickRate = tickRate
	player.ticker = ticker

	player.initPlayback()

	return player, nil
}

func (dp *DemoPlayer) initPlayback() {
	t := dp.ticker
	c := t.C

	go func() {
		for {
			<-c
			dp.result <- dp.nextTick()
		}
	}()
}

func isMapSupported(mapName string) bool {
	if _, ok := supportedMaps[mapName]; ok {
		return true
	}
	return false
}
