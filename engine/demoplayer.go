package engine

import (
	"fmt"
	"io"

	"github.com/markus-wa/demoinfocs-golang/msg"
	"github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
)

type DemoPlayer struct {
	mapName string
	parser  demoinfocs.Parser
	e       *engine
}

var player *DemoPlayer = nil

func (dp *DemoPlayer) Close() {
	if dp != nil && dp.parser != nil {
		dp.parser.Close()
	}
	player = nil
}

func (dp *DemoPlayer) NextTick() *StateResult {
	if dp.parser.Progress() == 1 {
		return nil
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

	player = &DemoPlayer{
		mapName: mapName,
		parser:  p,
		e:       &e,
	}

	return player, nil
}
