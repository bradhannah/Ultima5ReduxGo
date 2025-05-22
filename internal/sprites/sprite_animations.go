package sprites

import (
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
)

const (
	msPerFrameForObjects         = 125
	msPerFrameForNPCs            = 175
	shortNumberOfAnimationFrames = 2
	smallestNPCAnimationIndex    = indexes.AvatarSittingAndEatingFacingDown
)

func GetSpriteIndexWithAnimationBySpriteIndex(spriteIndex indexes.SpriteIndex, posHash int32) indexes.SpriteIndex {
	if spriteIndex >= indexes.Waterfall_KeyIndex && spriteIndex < indexes.Waterfall_KeyIndex+indexes.StandardNumberOfAnimationFrames {
		spriteIndex = indexes.Waterfall_KeyIndex
	}

	if spriteIndex == indexes.Waterfall_KeyIndex || spriteIndex == indexes.Fountain_KeyIndex {
		interval := time.Now().UnixMilli() / msPerFrameForObjects
		currentRotation := int(interval+int64(posHash)) % indexes.StandardNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)
	} else if spriteIndex == indexes.Clock1 || spriteIndex == indexes.Clock2 {
		interval := time.Now().UnixMilli() / (msPerFrameForObjects * 2)
		currentRotation := int(interval+int64(posHash)) % shortNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)
	} else if spriteIndex >= smallestNPCAnimationIndex {
		// if it is already animated, then we don't try to re-animate it because it can cause
		// visual glitches
		if spriteIndex%indexes.StandardNumberOfAnimationFrames != 0 {
			return spriteIndex
		}
		interval := time.Now().UnixMilli() / msPerFrameForNPCs
		currentRotation := int(interval+int64(posHash)) % indexes.StandardNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)
	}
	return spriteIndex
}
