package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionKlimbSmallMap(direction references.Direction) bool {
	// Check current position for ladder/grate first (player stands on these)
	currentTile := g.MapState.LayeredMaps.GetTileRefByPosition(
		references.SmallMapType,
		map_state.MapLayer,
		&g.MapState.PlayerLocation.Position,
		g.MapState.PlayerLocation.Floor)

	switch currentTile.Index {
	case indexes.AvatarOnLadderDown, indexes.LadderDown, indexes.Grate:
		locationRef := g.GameReferences.LocationReferences.GetLocationReference(g.MapState.PlayerLocation.Location)
		if locationRef.CanGoDownOneFloor(g.MapState.PlayerLocation.Floor) {
			g.MapState.PlayerLocation.Floor--
			g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
			g.SystemCallbacks.Message.AddRowStr("Klimb-Down!")
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		} else {
			g.SystemCallbacks.Message.AddRowStr("Can't go lower!")
			return false
		}

	case indexes.AvatarOnLadderUp, indexes.LadderUp:
		locationRef := g.GameReferences.LocationReferences.GetLocationReference(g.MapState.PlayerLocation.Location)
		if locationRef.CanGoUpOneFloor(g.MapState.PlayerLocation.Floor) {
			g.MapState.PlayerLocation.Floor++
			g.UpdateSmallMap(g.GameReferences.TileReferences, g.GameReferences.LocationReferences)
			g.SystemCallbacks.Message.AddRowStr("Klimb-Up!")
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		} else {
			g.SystemCallbacks.Message.AddRowStr("Can't go higher!")
			return false
		}
	}

	// Check target position for fence climbing
	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	targetTile := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor).GetTileTopMapOnlyTile(newPosition)
	if targetTile.Index == indexes.FenceHoriz || targetTile.Index == indexes.FenceVert {
		g.MapState.PlayerLocation.Position = *newPosition
		g.SystemCallbacks.Message.AddRowStr("Klimb!")
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true
	}

	g.SystemCallbacks.Message.AddRowStr("What?")
	return false
}

func (g *GameState) ActionKlimbLargeMap(direction references.Direction) bool {
	// Large map Klimb = mountain climbing with grapple (overworld only)

	// Check if player has grapple
	grappleQuantity := g.PartyState.Inventory.SpecialItems.GetQuantity(references.Grapple)
	if grappleQuantity == nil || !(*grappleQuantity).HasSome() {
		g.SystemCallbacks.Message.AddRowStr("With what?")
		return false
	}

	// Check if player is on foot (no vehicles)
	// TODO: Implement vehicle check when vehicle system is available
	// For now, assume player is always on foot on large map

	// Get target position and tile
	newPosition := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	targetTile := g.MapState.LayeredMaps.GetLayeredMap(references.LargeMapType, g.MapState.PlayerLocation.Floor).GetTopTile(newPosition)

	// Validate target is climbable mountains
	if targetTile.Index != indexes.SmallMountains {
		g.SystemCallbacks.Message.AddRowStr("Not climbable!")
		return false
	}

	// Check for impassable peaks
	if targetTile.Index == indexes.Peaks {
		g.SystemCallbacks.Message.AddRowStr("Impassable!")
		return false
	}

	// Perform dexterity checks for all living party members
	// Each member risks 1-5 damage on failure
	for i := 0; i < len(g.PartyState.Characters); i++ {
		member := &g.PartyState.Characters[i]
		if member.CurrentHp <= 0 {
			continue // Skip dead members
		}

		// Dexterity check: member dex vs random(1, 30)
		// TODO: Use proper random number generator when available
		// For now, use a simple check - assume 50% chance of failure for testing
		memberDex := int(member.Dexterity)
		randomRoll := 15 // Placeholder - should be random(1, 30)

		if memberDex < randomRoll {
			g.SystemCallbacks.Message.AddRowStr("Fell!")
			// Apply 1-5 damage
			damage := uint16(3) // Placeholder - should be random(1, 5)
			if member.CurrentHp > damage {
				member.CurrentHp -= damage
			} else {
				member.CurrentHp = 0
			}
		}
	}

	// Move player to the mountain position
	g.MapState.PlayerLocation.Position = *newPosition
	g.SystemCallbacks.Message.AddRowStr("Klimb!")
	g.SystemCallbacks.Flow.AdvanceTime(1)

	return true
}

func (g *GameState) ActionKlimbCombatMap(direction references.Direction) bool {
	// Combat map Klimb - handles ladders for level transitions
	// Uses current position tile examination like original implementation

	currentTile := g.MapState.LayeredMaps.GetTileRefByPosition(
		references.CombatMapType,
		map_state.MapLayer,
		&g.MapState.PlayerLocation.Position,
		g.MapState.PlayerLocation.Floor)

	switch currentTile.Index {
	case indexes.LadderUp:
		// Simple up ladder
		if g.canGoUpOneFloor() {
			g.goUpOneFloor()
			g.SystemCallbacks.Message.AddRowStr("Klimb-Up!")
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		} else {
			g.SystemCallbacks.Message.AddRowStr("Can't go higher!")
			return false
		}

	case indexes.LadderDown:
		// Simple down ladder
		if g.canGoDownOneFloor() {
			g.goDownOneFloor()
			g.SystemCallbacks.Message.AddRowStr("Klimb-Down!")
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		} else {
			g.SystemCallbacks.Message.AddRowStr("Can't go lower!")
			return false
		}

	// TODO: Handle bidirectional ladders when tile types are defined
	// This would require prompting for U/D direction choice
	// case indexes.BidirectionalLadder:
	//     // Would need to prompt "Klimb-U/D-" and get direction choice
	//     // Set exitdir = 5 (up) or 6 (down) for combat system

	default:
		g.SystemCallbacks.Message.AddRowStr("Nowhere to klimbe.")
		return false
	}
}

// Helper methods for floor transitions
func (g *GameState) canGoUpOneFloor() bool {
	// TODO: Implement floor validation logic
	// For now, assume you can always go up/down one floor in combat
	return true
}

func (g *GameState) canGoDownOneFloor() bool {
	// TODO: Implement floor validation logic
	// For now, assume you can always go up/down one floor in combat
	return true
}

func (g *GameState) goUpOneFloor() {
	g.MapState.PlayerLocation.Floor++
	// TODO: Update combat map when combat system is available
}

func (g *GameState) goDownOneFloor() {
	g.MapState.PlayerLocation.Floor--
	// TODO: Update combat map when combat system is available
}
