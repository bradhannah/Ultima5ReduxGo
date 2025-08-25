package game_state

import (
	"fmt"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameState) ActionGetSmallMap(direction references.Direction) bool {
	getThingPos := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	getThingTile := g.MapState.LayeredMaps.GetTileTopMapOnlyTileByPosition(references.SmallMapType, getThingPos, g.MapState.PlayerLocation.Floor)
	mapLayers := g.MapState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.MapState.PlayerLocation.Floor)

	// Check for item stacks first
	if g.ItemStacksMap.HasItemStackAtPosition(getThingPos) {
		item := g.ItemStacksMap.Pop(getThingPos)
		g.PartyState.Inventory.PutItemInInventory(item)

		itemRef := g.GameReferences.InventoryItemReferences.GetReferenceByItem(item.Item)
		g.SystemCallbacks.Message.AddRowStr(fmt.Sprintf("%s!", itemRef.ItemName))
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true
	}

	// Handle specific tile types
	switch getThingTile.Index {
	case indexes.WheatInField:
		g.SystemCallbacks.Message.AddRowStr("Crops picked! Those aren't yours Avatar!")
		mapLayers.SetTileByLayer(map_state.MapLayer, getThingPos, indexes.PlowedField)
		g.PartyState.Karma.DecreaseKarma(1)
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true

	case indexes.RightSconce, indexes.LeftScone:
		g.SystemCallbacks.Message.AddRowStr("Borrowed!")
		g.PartyState.Inventory.Provisions.Torches.IncrementByOne()
		mapLayers.SetTileByLayer(map_state.MapLayer, getThingPos, indexes.BrickFloor)
		g.SystemCallbacks.Flow.AdvanceTime(1)
		return true

	case indexes.TableFoodBoth, indexes.TableFoodBottom, indexes.TableFoodTop:
		if g.getFoodFromTable(direction, getThingPos, getThingTile, mapLayers) {
			g.SystemCallbacks.Message.AddRowStr("Mmmmm...! But that food isn't yours!")
			g.PartyState.Inventory.Provisions.Food.IncrementByOne()
			g.PartyState.Karma.DecreaseKarma(1)
			g.SystemCallbacks.Flow.AdvanceTime(1)
			return true
		}

	case indexes.Carpet2_MagicCarpet:
		// TODO: Implement magic carpet getting logic
		return false
	}

	// Nothing to get
	g.SystemCallbacks.Message.AddRowStr("Nothing there!")
	return false
}

// Helper function moved from GameScene to GameState
func (g *GameState) getFoodFromTable(direction references.Direction, getThingPos *references.Position, getThingTile *references.Tile, mapLayers *map_state.LayeredMap) bool {
	var newTileIndex indexes.SpriteIndex

	switch direction {
	case references.Down:
		if getThingTile.Index == indexes.TableFoodBoth {
			newTileIndex = indexes.TableFoodBottom
		} else if getThingTile.Index == indexes.TableFoodTop {
			newTileIndex = indexes.TableMiddle
		} else {
			return false
		}
	case references.Up:
		if getThingTile.Index == indexes.TableFoodBoth {
			newTileIndex = indexes.TableFoodTop
		} else if getThingTile.Index == indexes.TableFoodBottom {
			newTileIndex = indexes.TableMiddle
		} else {
			return false
		}
	default:
		return false
	}

	mapLayers.SetTileByLayer(map_state.MapLayer, getThingPos, newTileIndex)
	return true
}

func (g *GameState) ActionGetLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Get command - see Commands.md Get section
	// Large map variant of get command
	return true
}

func (g *GameState) ActionGetCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Get command - see Commands.md Get section
	// Combat map variant of get command
	return true
}

func (g *GameState) ActionGetDungeonMap(direction references.Direction) bool {
	// TODO: Implement dungeon map Get command - see Commands.md Get — Dungeon section
	// Should handle:
	// - Picking from underfoot opened chest
	// - Light requirement for seeing items
	// - Distinct from surface object pickup
	return true
}
