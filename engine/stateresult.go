// Contains types that represent State that is being transfered to the drawing function
package engine

import "time"

// Contains state for current tick
type StateResult struct {
	TeamT         Team
	TeamCt        Team
	RoundTimeLeft Second
	Nades         []Nade
	Infernos      []Inferno
	CurrentRound  int
	Bomb          Bomb
}

type Team struct {
	Players []PlayerData
	ClanTag string
	Score   int
}

type Nade struct {
	Positions []Position
	Type      string
}
type Bomb struct {
	CarrierId int // 0 if no carrier
	Position  Position
	isPlanted bool
	isDefused bool
}

// Position of an entity on the map
type Position struct {
	X float64
	Y float64
}
type Inferno struct {
	BorderPositions []Position
}

type Second time.Duration
type Millisecond time.Duration

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
	FlashTimeLeft     Millisecond
}
