package engine

import (
	"fmt"
	"os"

	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
)

type DemoPlayer struct {
	paused        bool
	currentRound  int
	playbackSpeed float64
	mapName       string
	parser        demoinfocs.Parser
}

var player *DemoPlayer = nil

func (p *DemoPlayer) Pause() {
	p.paused = true
}

func (p *DemoPlayer) Play() {
	p.paused = false
}

func (p *DemoPlayer) Close() {
	if player != nil && player.parser != nil {
		p.parser.Close()
	}
	player = nil
}

func NextTick() {

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

	p.GameState().Participants().Playing()

	if player == nil {
		player = &DemoPlayer{
			paused:        true,
			currentRound:  1,
			playbackSpeed: 1.0,
			mapName:       mapName,
			parser:        p,
		}
	}

	return player, nil
}
