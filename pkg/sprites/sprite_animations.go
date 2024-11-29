package sprites

import (
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

const (
	msPerFrameForObjects            = 125
	msPerFrameForNPCs               = 175
	standardNumberOfAnimationFrames = 4
	smallestNPCAnimationIndex       = indexes.AvatarSittingAndEatingFacingDown
)

func GetSpriteIndexWithAnimationBySpriteIndex(spriteIndex indexes.SpriteIndex, posHash int32) indexes.SpriteIndex {
	if spriteIndex >= indexes.Waterfall_KeyIndex && spriteIndex < indexes.Waterfall_KeyIndex+4 {
		spriteIndex = indexes.Waterfall_KeyIndex
	}

	if spriteIndex == indexes.Waterfall_KeyIndex || spriteIndex == indexes.Fountain_KeyIndex {
		interval := time.Now().UnixMilli() / msPerFrameForObjects
		currentRotation := int(interval+int64(posHash)) % standardNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)
	} else if spriteIndex >= smallestNPCAnimationIndex {
		// if it is already animated, then we don't try to re-animate it because it can cause
		// visual glitches
		if spriteIndex%4 != 0 {
			return spriteIndex
		}
		interval := time.Now().UnixMilli() / msPerFrameForNPCs
		currentRotation := int(interval+int64(posHash)) % standardNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)
	}
	return spriteIndex
}
