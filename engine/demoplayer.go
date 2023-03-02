package engine

import (
	"encoding/json"
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

func (dp *DemoPlayer) NextTick() string {
	if dp.parser.Progress() == 1 {
		return ""
	}
	dp.parser.ParseNextFrame()
	res := dp.e.getUsefulState(dp.parser.GameState())

	json, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	return string(json)
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

func getMapMetadata(p demoinfocs.Parser, mapName string) Map {
	var metadata Map
	p.RegisterNetMessageHandler(func(msg *msg.CSVCMsg_ServerInfo) {
		fmt.Println(mapName, msg.GetMapCrc())
		// Get metadata for the map that the game was played on for coordinate translations
		metadata = GetMapMetadata(mapName, msg.GetMapCrc())
		fmt.Println(metadata)
	})
	return metadata
}

type promise struct {
	done chan struct{}
	data Map
}
