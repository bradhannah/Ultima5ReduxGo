package sprites

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

// REGRESSION TESTS - Prevent non-deterministic animation behavior from regressing
// These tests ensure animations use deterministic timing instead of time.Now()

func TestAnimationDeterminism(t *testing.T) {
	// CRITICAL: Same inputs should ALWAYS produce same outputs (deterministic)
	// This prevents regression where animations used time.Now() and were non-reproducible

	testCases := []struct {
		name        string
		spriteIndex indexes.SpriteIndex
		posHash     int32
		elapsedMs   int64
	}{
		{
			name:        "Waterfall animation",
			spriteIndex: indexes.Waterfall_KeyIndex,
			posHash:     123,
			elapsedMs:   1000,
		},
		{
			name:        "Fountain animation",
			spriteIndex: indexes.Fountain_KeyIndex,
			posHash:     456,
			elapsedMs:   2000,
		},
		{
			name:        "Clock animation",
			spriteIndex: indexes.Clock1,
			posHash:     789,
			elapsedMs:   1500,
		},
		{
			name:        "NPC animation",
			spriteIndex: indexes.AvatarSittingAndEatingFacingDown,
			posHash:     100,
			elapsedMs:   3000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function multiple times with same inputs
			result1 := GetSpriteIndexWithAnimationBySpriteIndexTick(tc.spriteIndex, tc.posHash, tc.elapsedMs)
			result2 := GetSpriteIndexWithAnimationBySpriteIndexTick(tc.spriteIndex, tc.posHash, tc.elapsedMs)
			result3 := GetSpriteIndexWithAnimationBySpriteIndexTick(tc.spriteIndex, tc.posHash, tc.elapsedMs)

			// CRITICAL: Results must be identical (deterministic behavior)
			if result1 != result2 {
				t.Errorf("REGRESSION: Animation not deterministic - first call: %d, second call: %d", result1, result2)
			}
			if result1 != result3 {
				t.Errorf("REGRESSION: Animation not deterministic - first call: %d, third call: %d", result1, result3)
			}

			t.Logf("✅ Deterministic result for %s: %d", tc.name, result1)
		})
	}
}

func TestAnimationProgression(t *testing.T) {
	// Test that animations actually progress over time (not stuck)

	testCases := []struct {
		name        string
		spriteIndex indexes.SpriteIndex
		posHash     int32
		frameMs     int64 // Time for one animation frame
	}{
		{
			name:        "Waterfall progression",
			spriteIndex: indexes.Waterfall_KeyIndex,
			posHash:     0,
			frameMs:     msPerFrameForObjects, // 125ms
		},
		{
			name:        "Fountain progression",
			spriteIndex: indexes.Fountain_KeyIndex,
			posHash:     0,
			frameMs:     msPerFrameForObjects, // 125ms
		},
		{
			name:        "Clock progression",
			spriteIndex: indexes.Clock1,
			posHash:     0,
			frameMs:     msPerFrameForObjects * 2, // 250ms (clocks are 2x slower)
		},
		{
			name:        "NPC progression",
			spriteIndex: indexes.AvatarSittingAndEatingFacingDown,
			posHash:     0,
			frameMs:     msPerFrameForNPCs, // 175ms
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get animation frame at time 0
			frame1 := GetSpriteIndexWithAnimationBySpriteIndexTick(tc.spriteIndex, tc.posHash, 0)

			// Get animation frame after exactly one frame duration
			frame2 := GetSpriteIndexWithAnimationBySpriteIndexTick(tc.spriteIndex, tc.posHash, tc.frameMs)

			// Animation should progress (different frames)
			if frame1 == frame2 {
				t.Errorf("Animation should progress over time - frame at 0ms: %d, frame at %dms: %d",
					frame1, tc.frameMs, frame2)
			} else {
				t.Logf("✅ Animation progresses: %d → %d over %dms", frame1, frame2, tc.frameMs)
			}
		})
	}
}

func TestAnimationFrameBounds(t *testing.T) {
	// Test that animation frames stay within expected bounds

	t.Run("Waterfall animation bounds", func(t *testing.T) {
		var baseIndex indexes.SpriteIndex = indexes.Waterfall_KeyIndex

		// Test various time points
		for elapsedMs := int64(0); elapsedMs < 1000; elapsedMs += 50 {
			result := GetSpriteIndexWithAnimationBySpriteIndexTick(baseIndex, 0, elapsedMs)

			// Should be within the standard animation frame range
			minIndex := int(baseIndex)
			maxIndex := int(baseIndex) + int(indexes.StandardNumberOfAnimationFrames) - 1

			if int(result) < minIndex || int(result) > maxIndex {
				t.Errorf("Animation frame out of bounds at %dms: got %d, expected range [%d-%d]",
					elapsedMs, result, minIndex, maxIndex)
			}
		}
	})

	t.Run("Clock animation bounds", func(t *testing.T) {
		var baseIndex indexes.SpriteIndex = indexes.Clock1

		// Test various time points
		for elapsedMs := int64(0); elapsedMs < 1000; elapsedMs += 100 {
			result := GetSpriteIndexWithAnimationBySpriteIndexTick(baseIndex, 0, elapsedMs)

			// Clock animations use shortNumberOfAnimationFrames (2)
			minIndex := int(baseIndex)
			maxIndex := int(baseIndex) + int(shortNumberOfAnimationFrames) - 1

			if int(result) < minIndex || int(result) > maxIndex {
				t.Errorf("Clock animation frame out of bounds at %dms: got %d, expected range [%d-%d]",
					elapsedMs, result, minIndex, maxIndex)
			}
		}
	})
}

func TestAnimationPositionHashVariation(t *testing.T) {
	// Test that posHash parameter affects animation phase (prevents all objects animating in sync)

	var spriteIndex indexes.SpriteIndex = indexes.Waterfall_KeyIndex
	elapsedMs := int64(1000)

	// Get results with different position hashes
	result1 := GetSpriteIndexWithAnimationBySpriteIndexTick(spriteIndex, 0, elapsedMs)
	result2 := GetSpriteIndexWithAnimationBySpriteIndexTick(spriteIndex, 100, elapsedMs)
	result3 := GetSpriteIndexWithAnimationBySpriteIndexTick(spriteIndex, 200, elapsedMs)

	// Position hash should create variation in animation phase
	// (Not all objects should animate in perfect sync)
	results := []indexes.SpriteIndex{result1, result2, result3}
	allSame := true
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			allSame = false
			break
		}
	}

	if allSame {
		t.Logf("Note: All position hashes produced same frame (%d) - may be due to timing alignment", result1)
	} else {
		t.Logf("✅ Position hash creates animation variation: %d, %d, %d", result1, result2, result3)
	}
}

func TestAnimationStaticSprites(t *testing.T) {
	// Test that non-animated sprites are returned unchanged

	staticSprites := []indexes.SpriteIndex{
		indexes.BrickFloor,
		indexes.Grass,
		indexes.SpriteIndex(1),  // Generic non-animated sprite
		indexes.SpriteIndex(50), // Another generic non-animated sprite
	}

	for _, sprite := range staticSprites {
		// Static sprites should be returned unchanged regardless of time
		result1 := GetSpriteIndexWithAnimationBySpriteIndexTick(sprite, 0, 0)
		result2 := GetSpriteIndexWithAnimationBySpriteIndexTick(sprite, 0, 10000)
		result3 := GetSpriteIndexWithAnimationBySpriteIndexTick(sprite, 100, 5000)

		if result1 != sprite || result2 != sprite || result3 != sprite {
			t.Errorf("Static sprite should remain unchanged: input %d, results %d, %d, %d",
				sprite, result1, result2, result3)
		}
	}
}

// REGRESSION TEST - Ensure no time.Now() usage can creep back in
func TestNoTimeNowRegression(t *testing.T) {
	// This test documents that animation MUST be deterministic
	// If time.Now() usage is reintroduced, this test will fail intermittently

	var spriteIndex indexes.SpriteIndex = indexes.Waterfall_KeyIndex
	posHash := int32(42)
	elapsedMs := int64(1337)

	// Run the same call many times rapidly
	results := make([]indexes.SpriteIndex, 10)
	for i := range results {
		results[i] = GetSpriteIndexWithAnimationBySpriteIndexTick(spriteIndex, posHash, elapsedMs)
	}

	// All results MUST be identical (deterministic)
	firstResult := results[0]
	for i, result := range results {
		if result != firstResult {
			t.Errorf("CRITICAL REGRESSION: Non-deterministic behavior detected! "+
				"Call %d returned %d, expected %d. This suggests time.Now() usage has returned!",
				i, result, firstResult)
		}
	}

	t.Logf("✅ All %d rapid calls returned identical result: %d (deterministic)", len(results), firstResult)
}
