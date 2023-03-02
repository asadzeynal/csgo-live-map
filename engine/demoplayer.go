package engine

import (
	"fmt"
	"os"

	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
)

type DemoPlayer struct {
	playbackSpeed float64
	mapName       string
	parser        demoinfocs.Parser
}

var player *DemoPlayer = nil

func (dp *DemoPlayer) Close() {
	if dp != nil && dp.parser != nil {
		dp.parser.Close()
	}
	player = nil
}

func (dp *DemoPlayer) NextTick() StateResult {
	dp.parser.ParseNextFrame()
	players := dp.parser.GameState().Participants().Playing()
	return processOneTick(IncomingState{Players: players})
}

func GetPlayer(file *os.File) (*DemoPlayer, error) {
	p := demoinfocs.NewParser(file)
	header, err := p.ParseHeader()
	if err != nil {
		return nil, fmt.Errorf("unable to parse demo headers: %v", err)
	}

	mapName := header.MapName
	if mapName != "de_ancient" {
		return nil, fmt.Errorf("only de_ancient is supported now")
	}

	for !p.GameState().IsMatchStarted() {
		p.ParseNextFrame()
	}

	if player == nil {
		player = &DemoPlayer{
			playbackSpeed: 1.0,
			mapName:       mapName,
			parser:        p,
		}
	}

	return player, nil
}
