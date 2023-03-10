package engine

import (
	"bytes"
	"fmt"
	"time"

	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
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
	file          []byte
	IsPaused      bool
	MapName       string
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

func (dp *DemoPlayer) NextRound() {
	dp.Pause()
	nextRound := dp.e.currentRound + 1
	for dp.e.currentRound != nextRound {
		dp.nextTick()
	}
	dp.result <- dp.nextTick()
}

func (dp *DemoPlayer) PrevRound() {
	prevRound := dp.e.currentRound - 1
	dp.stopInternal()
	if prevRound < 1 {
		return
	}
	for dp.e.currentRound != prevRound {
		dp.parser.ParseNextFrame()
	}
	dp.result <- dp.nextTick()
}

func (dp *DemoPlayer) Pause() {
	if dp.IsPaused {
		fmt.Println("isPaused")
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

func (dp *DemoPlayer) stopInternal() {
	dp.Pause()
	dp.parser.Cancel()
	dp.parser.Close()
	dp.e.currentRound = 1
	dp.e.scoreCt = 0
	dp.e.scoreT = 0
	dp.parser = demoinfocs.NewParser(bytes.NewReader(dp.file))
	dp.registerEventHandlers()
	for !dp.parser.GameState().IsMatchStarted() {
		dp.parser.ParseNextFrame()
	}
}

func (dp *DemoPlayer) Stop() {
	dp.stopInternal()
	dp.result <- dp.nextTick()
}

func (dp *DemoPlayer) refreshTicker() {
	dp.ticker.Reset(dp.calculateTickerDuration())
}

func (dp *DemoPlayer) ChangeSpeed(speed float64) {
	dp.playbackSpeed = speed
	dp.refreshTicker()
	dp.Pause()
}

func (dp *DemoPlayer) nextTick() StateResult {
	if dp.parser.Progress() == 1 {
		return StateResult{}
	}
	dp.parser.ParseNextFrame()
	res := dp.e.getUsefulState(dp.parser.GameState(), dp.parser.CurrentTime())

	return res
}

func (dp *DemoPlayer) WaitForStateUpdate() StateResult {
	return <-dp.result
}

func (dp *DemoPlayer) registerEventHandlers() {
	dp.parser.RegisterEventHandler(func(event events.RoundFreezetimeEnd) {
		dp.e.roundFreezeTimeEndAt = dp.parser.CurrentTime()
	})

	dp.parser.RegisterEventHandler(func(event events.RoundEnd) {
		dp.e.roundEndedAt = dp.parser.CurrentTime()
		if event.Winner == 2 {
			dp.e.scoreT++
		} else {
			dp.e.scoreCt++
		}
	})
	dp.parser.RegisterEventHandler(func(event events.RoundEndOfficial) {
		dp.e.currentRound++
	})

}

func GetPlayer(fileRaw []byte) (*DemoPlayer, error) {
	if player != nil {
		return player, nil
	}
	file := bytes.NewReader(fileRaw)

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
		file:          fileRaw,
		MapName:       mapName,
		parser:        p,
		IsPaused:      true,
		playbackSpeed: 1.0,
		result:        make(chan StateResult),
	}

	e := engine{
		currentRound: 1,
	}

	p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_ServerInfo) {
		mmd := GetMapMetadata(mapName, msg.GetMapCrc())
		e.mapMetadata = &mmd
	})

	player.registerEventHandlers()

	for !p.GameState().IsMatchStarted() {
		p.ParseNextFrame()
	}
	tickRate := p.TickRate()

	player.e = &e
	player.tickRate = tickRate

	ticker := time.NewTicker(player.calculateTickerDuration())
	ticker.Stop()
	player.ticker = ticker

	player.initPlayback()

	return player, nil
}

func (dp *DemoPlayer) calculateTickerDuration() time.Duration {
	timePerFrame := int(time.Second) / int(player.tickRate)
	withSpeed := float64(timePerFrame) / dp.playbackSpeed
	return time.Duration(withSpeed)
}

func (dp *DemoPlayer) initPlayback() {
	t := dp.ticker
	c := t.C

	go func() {
		dp.result <- dp.nextTick()
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
