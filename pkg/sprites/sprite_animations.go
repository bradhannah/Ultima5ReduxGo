package sprites

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"time"
)

const msPerFrame = 75

func GetTileNumberWithAnimationByTile(nTile int) int {
	//mainTile := int(s.rawData[nFloor][position.X][position.Y])
	if (nTile >= indexes.Waterfall_KeyIndex && nTile <= indexes.Waterfall_KeyIndex+3) || nTile == indexes.Fountain_KeyIndex || nTile >= 308 {
		// animation
		// msPerFrame
		interval := time.Now().UnixMilli() / msPerFrame
		currentRotation := int(interval) % 4
		return nTile + currentRotation
	}
	return nTile
}
