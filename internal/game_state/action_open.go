package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionOpenSmallMap(direction references.Direction) bool {
	openThingPos := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	openThingTile := g.GetLayeredMapByCurrentLocation().GetTileTopMapOnlyTile(openThingPos)

	// Handle doors first
	if openThingTile.Index.IsDoor() {
		switch g.MapState.OpenDoor(direction) {
		case map_state.OpenDoorNotADoor:
			g.SystemCallbacks.Message.AddRowStr("Bang to open!")
			return false
		case map_state.OpenDoorLocked:
			g.SystemCallbacks.Message.AddRowStr("Locked!")
			return false
		case map_state.OpenDoorLockedMagical:
			g.SystemCallbacks.Message.AddRowStr("Magically Locked!")
			return false
		case map_state.OpenDoorOpened:
			g.SystemCallbacks.Message.AddRowStr("Opened!")
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		}
		return false
	}

	// Handle chests
	openThingTopTile := g.GetLayeredMapByCurrentLocation().GetTopTile(openThingPos)
	if openThingTopTile.Index == indexes.Chest {
		// Special case: Lord British's castle treasure
		if g.MapState.PlayerLocation.Location == references.Lord_Britishs_Castle && g.MapState.PlayerLocation.Floor == references.Basement {
			itemStack := references.CreateNewItemStack(references.LordBritishTreasure)
			g.SystemCallbacks.Message.AddRowStr("Found:")
			g.SystemCallbacks.Message.AddRowStr(g.GameReferences.InventoryItemReferences.GetListOfItems(&itemStack))
			g.ItemStacksMap.Push(openThingPos, &itemStack)
			g.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(*openThingPos)
			g.CurrentNPCAIController.FreshenExistingNPCsOnMap()
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		}
		// TODO: Handle other chests when general chest opening is implemented
		g.SystemCallbacks.Message.AddRowStr("Nothing there!")
		return false
	}

	// Nothing openable
	g.SystemCallbacks.Message.AddRowStr("Bang to open!")
	return false
}

func (g *GameState) ActionOpenLargeMap(direction references.Direction) bool {
	targetPos := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)

	// Get the tile at the target position
	layeredMap := g.GetLayeredMapByCurrentLocation()
	topTile := layeredMap.GetTopTile(targetPos)

	// Check if there's an item stack at this position (like chests)
	if g.ItemStacksMap.HasItemStackAtPosition(targetPos) {
		// For now, just indicate there's something here but can't be opened on large map
		// This could be expanded later to handle specific cases
		return false
	}

	// Check if the tile is openable (like doors)
	if topTile.Index.IsDoor() {
		// Use the door opening logic from small map
		switch g.MapState.OpenDoor(direction) {
		case map_state.OpenDoorOpened:
			g.SystemCallbacks.Message.AddRowStr("Opened!")
			return true
		case map_state.OpenDoorLocked:
			g.SystemCallbacks.Message.AddRowStr("Locked!")
			return false
		case map_state.OpenDoorLockedMagical:
			g.SystemCallbacks.Message.AddRowStr("Magically Locked!")
			return false
		case map_state.OpenDoorNotADoor:
			// Not a door, continue to check other openable things
			break
		}
	}

	// On large maps, most openable things are handled by entering buildings
	// or specific interactions, so most open attempts will fail
	return false
}

func (g *GameState) ActionOpenCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Open command - see Commands.md Open section
	// Combat map variant of open command - likely limited to chests
	return true
}

func (g *GameState) ActionOpenDungeonMap(direction references.Direction) bool {
	// TODO: Implement dungeon map Open command - see Commands.md Open — Dungeon section
	// Should handle:
	// - Dungeon doors and underfoot chests
	// - Light requirement for seeing targets
	// - Integration with spells (An Sanct/In Ex Por)
	return true
}
