package sprites

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

const (
	msPerFrameForObjects         = 125
	msPerFrameForNPCs            = 175
	shortNumberOfAnimationFrames = 2
	smallestNPCAnimationIndex    = indexes.AvatarSittingAndEatingFacingDown
)

// GetSpriteIndexWithAnimationBySpriteIndexTick is a deterministic animation function.
// Use elapsedMs from your game clock so animation timing is in sync with gameplay timing.
func GetSpriteIndexWithAnimationBySpriteIndexTick(spriteIndex indexes.SpriteIndex, posHash int32, elapsedMs int64) indexes.SpriteIndex {
	if spriteIndex >= indexes.Waterfall_KeyIndex && spriteIndex < indexes.Waterfall_KeyIndex+indexes.StandardNumberOfAnimationFrames {
		spriteIndex = indexes.Waterfall_KeyIndex
	}

	switch {
	case spriteIndex == indexes.Waterfall_KeyIndex || spriteIndex == indexes.Fountain_KeyIndex:
		interval := elapsedMs / msPerFrameForObjects
		currentRotation := int(interval+int64(posHash)) % indexes.StandardNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)

	case spriteIndex == indexes.Clock1 || spriteIndex == indexes.Clock2:
		interval := elapsedMs / (msPerFrameForObjects * 2)
		currentRotation := int(interval+int64(posHash)) % shortNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)

	case spriteIndex >= smallestNPCAnimationIndex:
		// Avoid re-animating mid-cycle frames to prevent glitching.
		if spriteIndex%indexes.StandardNumberOfAnimationFrames != 0 {
			return spriteIndex
		}
		interval := elapsedMs / msPerFrameForNPCs
		currentRotation := int(interval+int64(posHash)) % indexes.StandardNumberOfAnimationFrames
		return indexes.SpriteIndex(int(spriteIndex) + currentRotation)
	}

	return spriteIndex
}
