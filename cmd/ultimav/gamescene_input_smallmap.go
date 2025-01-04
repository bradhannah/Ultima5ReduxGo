package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameScene) smallMapInputHandler(key ebiten.Key) {
	switch key {
	case ebiten.KeyEscape:
		g.DoEscapeMenu()
		return
	case ebiten.KeySpace:
		g.addRowStr("Pass")
		// g.gameState.FinishTurn()
	case ebiten.KeyBackquote:
		g.ToggleDebug()
		return
	case ebiten.KeyEnter:
		g.addRowStr("Enter")
	case ebiten.KeyUp:
		g.handleMovement(references.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(references.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(references.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(references.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyB:
		g.actionBoard()
	case ebiten.KeyK:
		g.smallMapKlimb()
	case ebiten.KeyL:
		g.addRowStr("Look-")
		g.secondaryKeyState = LookDirectionInput
	case ebiten.KeyX:
		g.gameState.ExitSmallMap()
	case ebiten.KeyG:
		// get the thing - direction
		g.addRowStr("Get-")
		g.secondaryKeyState = GetDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyE:
		g.addRowStr(fmt.Sprintf("Enter what?"))
	case ebiten.KeyP:
		g.addRowStr("Push-")
		g.secondaryKeyState = PushDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyO:
		g.debugConsole.Output.AddRowStr("Open")
		g.addRowStr("Open-")
		g.secondaryKeyState = OpenDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyJ:
		g.debugMessage = "Jimmy"
		g.addRowStr("Jimmy-")
		g.secondaryKeyState = JimmyDoorDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	default:
		return
	}

	// only process end of turn if the turn is actually done.
	if g.secondaryKeyState == PrimaryInput {
		g.gameState.FinishTurn()
	}
}

func (g *GameScene) smallMapHandleSecondaryInput() {
	arrowKey := getArrowKeyPressed()
	bIsArrowKeyPressed := arrowKey != nil

	switch g.secondaryKeyState {
	case JimmyDoorDirectionInput:
		if g.gameState.Inventory.Provisions.Keys <= 0 {
			g.addRowStr("No Keys!")
			g.secondaryKeyState = PrimaryInput
			g.keyboard.SetLastKeyPressedNow()
			return
		}

		if g.isDirectionKeyValidAndOutput() {
			g.smallMapJimmySecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case OpenDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapOpenSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case KlimbDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapKlimbSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case PushDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapPushSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case GetDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapGetSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case LookDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.commonMapLookSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	default:
		panic("unhandled default case")
	}

	if bIsArrowKeyPressed {
		g.keyboard.SetLastKeyPressedNow()
	}
}

func (g *GameScene) smallMapKlimb() {
	currentTile := g.gameState.LayeredMaps.GetTileRefByPosition(references.SmallMapType, game_state.MapLayer, &g.gameState.Position, g.gameState.Floor)

	switch currentTile.Index {
	case indexes.AvatarOnLadderDown, indexes.LadderDown, indexes.Grate:
		if g.GetCurrentLocationReference().CanGoDownOneFloor(g.gameState.Floor) {
			g.gameState.Floor--
			g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
			g.output.AddRowStr("Klimb-Down!")
			return
		} else {
			log.Fatal("Can't go lower my dude")
		}

	case indexes.AvatarOnLadderUp, indexes.LadderUp:
		if g.GetCurrentLocationReference().CanGoUpOneFloor(g.gameState.Floor) {
			g.gameState.Floor++
			g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
			g.output.AddRowStr("Klimb-Up!")
			return
		} else {
			log.Fatal("Can't go higher my dude")
		}
	}
	g.output.AddRowStr("Klimb-")
	g.secondaryKeyState = KlimbDirectionInput
}

func (g *GameScene) smallMapKlimbSecondary(direction references.Direction) {
	if !g.gameState.ActionKlimbSmallMap(direction) {
		g.output.AddRowStr("What?")
	}
}

func (g *GameScene) smallMapPushSecondary(direction references.Direction) {
	pushThingPos := direction.GetNewPositionInDirection(&g.gameState.Position)

	// pushThingTile := g.gameState.LayeredMaps.GetTileTopMapOnlyTileByPosition(references.SmallMapType, pushThingPos, g.gameState.Floor)
	pushThingTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(pushThingPos)

	if !pushThingTile.IsPushable {
		g.output.AddRowStr("Won't budge!")
		return
	}

	if g.gameState.ActionPushSmallMap(direction) {
		// moved
	} else {
		// didn't move
	}
}

func (g *GameScene) smallMapOpenSecondary(direction references.Direction) {
	openThingPos := direction.GetNewPositionInDirection(&g.gameState.Position)
	openThingTile := g.gameState.GetLayeredMapByCurrentLocation().GetTileTopMapOnlyTile(openThingPos)

	if openThingTile.Index.IsDoor() {
		switch g.gameState.OpenDoor(direction) {
		case game_state.OpenDoorNotADoor:
			g.addRowStr("Bang to open!")
		case game_state.OpenDoorLocked:
			g.addRowStr("Locked!")
		case game_state.OpenDoorLockedMagical:
			g.addRowStr("Magically Locked!")
		case game_state.OpenDoorOpened:
			g.addRowStr("Opened!")
		default:
			log.Fatal("Unrecognized door open state")
		}
		return
	}

	openThingTopTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(openThingPos)
	if openThingTopTile.Index == indexes.Chest {
		if g.gameState.Location == references.Lord_Britishs_Castle && g.gameState.Floor == references.Basement {
			itemStack := references.CreateNewItemStack(references.LordBritishTreasure)
			g.output.AddRowStr("Found:")
			g.output.AddRowStr(g.gameReferences.InventoryItemReferences.GetListOfItems(&itemStack))
			g.gameState.ItemStacksMap.Push(openThingPos, &itemStack)
			g.gameState.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(*openThingPos)
			g.gameState.CurrentNPCAIController.FreshenExistingNPCsOnMap()
		}
	}
}

// func (g *GameScene) smallMapBoard() bool {
// 	getThingTile := g.gameState.LayeredMaps.GetTileRefByPosition(references.SmallMapType, game_state.MapUnitLayer, &g.gameState.Position, g.gameState.Floor)
//
// 	if !g.gameState.BoardVehicle(references.GetVehicleFromSpriteIndex(getThingTile.Index)) {
// 		g.output.AddRowStr("Board what?")
// 		return false
// 	}
//
// 	g.gameState.CurrentNPCAIController.RemoveNPCAtPosition(g.gameState.Position)
//
// 	switch getThingTile.Index {
// 	case indexes.Carpet2_MagicCarpet:
// 		g.output.AddRowStr("carpet")
// 	case indexes.HorseRight, indexes.HorseLeft:
// 		g.output.AddRowStr("horse")
// 	case indexes.FrigateDownFurled, indexes.FrigateUpFurled, indexes.FrigateLeftFurled, indexes.FrigateRightFurled:
// 		g.output.AddRowStr("Frigate")
// 	case indexes.SkiffLeft, indexes.SkiffRight, indexes.SkiffUp, indexes.SkiffDown:
// 		g.output.AddRowStr("Skiff")
// 	default:
// 	}
//
// 	return true
// }

func (g *GameScene) smallMapJimmySecondary(direction references.Direction) {
	jimmyResult := g.gameState.JimmyDoor(direction, &g.gameState.Characters[0])

	switch jimmyResult {
	case game_state.JimmyUnlocked:
		g.addRowStr("Unlocked!")
	case game_state.JimmyNotADoor:
		g.addRowStr("Not lock!")
	case game_state.JimmyBrokenPick, game_state.JimmyLockedMagical:
		g.addRowStr("Key broke!")
	default:
		panic("unhandled default case")
	}
}

func (g *GameScene) smallMapGetSecondary(direction references.Direction) {
	getThingPos := direction.GetNewPositionInDirection(&g.gameState.Position)
	getThingTile := g.gameState.LayeredMaps.GetTileTopMapOnlyTileByPosition(references.SmallMapType, getThingPos, g.gameState.Floor)
	mapLayers := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor)

	if g.gameState.ItemStacksMap.HasItemStackAtPosition(getThingPos) {
		item := g.gameState.ItemStacksMap.Pop(getThingPos)
		g.gameState.Inventory.PutItemInInventory(item)

		itemRef := g.gameReferences.InventoryItemReferences.GetReferenceByItem(item.Item)
		// if item.
		g.addRowStr(fmt.Sprintf("%s!", itemRef.ItemName))
		return
	}

	switch getThingTile.Index {
	case indexes.WheatInField:
		g.addRowStr("Crops picked! Those aren't yours Avatar!")
		mapLayers.SetTileByLayer(game_state.MapLayer, getThingPos, indexes.PlowedField)
		g.gameState.Karma = g.gameState.Karma.GetDecreasedKarma(1)
	case indexes.RightSconce, indexes.LeftScone:
		g.addRowStr("Borrowed!")
		g.gameState.Inventory.Provisions.Torches++
		mapLayers.SetTileByLayer(game_state.MapLayer, getThingPos, indexes.BrickFloor)
	case indexes.TableFoodBoth, indexes.TableFoodBottom, indexes.TableFoodTop:
		if g.getFoodFromTable(direction) {
			g.addRowStr("Mmmmm...! But that food isn't yours!")
			g.gameState.Inventory.Provisions.Food++
			g.gameState.Karma = g.gameState.Karma.GetDecreasedKarma(1)
		}
	case indexes.Carpet2_MagicCarpet:
	}
}

func (g *GameScene) getFoodFromTable(direction references.Direction) bool {
	getThingPos := direction.GetNewPositionInDirection(&g.gameState.Position)
	getThingTile := g.gameState.LayeredMaps.GetTileTopMapOnlyTileByPosition(references.SmallMapType, getThingPos, g.gameState.Floor)
	mapLayers := g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor)

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

	mapLayers.SetTileByLayer(game_state.MapLayer, getThingPos, newTileIndex)
	return true
}
