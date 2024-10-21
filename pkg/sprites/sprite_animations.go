package sprites

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"time"
)

const msPerFrame = 125

func GetSpriteIndexWithAnimationBySpriteIndex(spriteIndex indexes.SpriteIndex) indexes.SpriteIndex {
	//mainTile := int(s.rawData[nFloor][position.X][position.Y])
	if (spriteIndex >= indexes.Waterfall_KeyIndex && spriteIndex <= indexes.Waterfall_KeyIndex+3) || spriteIndex == indexes.Fountain_KeyIndex || spriteIndex >= 308 {
		// animation
		// msPerFrame
		interval := time.Now().UnixMilli() / msPerFrame
		currentRotation := int(interval) % 4
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)
	}
	return spriteIndex
}
