package engine

import "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/common"

type CurrentGameState struct {
	Players []*common.Player
}

func processOneTick(state CurrentGameState) {

}

func Run(update chan CurrentGameState) {
	for {
		currentState := <-update
		processOneTick(currentState)
	}
}
