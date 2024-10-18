package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func (g *GameState) DebugMoveOnMap(position references.Position) {
	g.Position = position
}
