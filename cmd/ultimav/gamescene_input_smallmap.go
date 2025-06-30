package main

import (
	"fmt"
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/hajimehoshi/ebiten/v2"

	gamestate "github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

func (g *GameScene) smallMapInputHandler(key ebiten.Key) {
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			g.gameState.DebugQuickExitSmallMap()
			return
		}
	}

	switch key {
	case ebiten.KeyEscape:
		g.DoEscapeMenu()
		return
	case ebiten.KeySpace:
		g.addRowStr("Pass")
		// g.gameState.FinishTurn()
	case ebiten.KeyBackquote:
		g.toggleDebug()
		return
	case ebiten.KeyEnter:
		g.addRowStr("Enter")
	case ebiten.KeyUp:
		g.handleMovement(references2.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(references2.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(references2.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(references2.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyB:
		g.actionBoard()
	case ebiten.KeyK:
		g.smallMapKlimb()
	case ebiten.KeyL:
		g.addRowStr("Look-")
		g.secondaryKeyState = LookDirectionInput
	case ebiten.KeyX:
		g.actionExit()
	case ebiten.KeyG:
		// get the thing - direction
		g.addRowStr("Get-")
		g.secondaryKeyState = GetDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyE:
		g.addRowStr("Enter what?")
	case ebiten.KeyP:
		g.addRowStr("Push-")
		g.secondaryKeyState = PushDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyO:
		g.debugConsole.Output.AddRowStrWithTrim("Open")
		g.addRowStr("Open-")
		g.secondaryKeyState = OpenDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyJ:
		g.debugMessage = "Jimmy"
		g.addRowStr("Jimmy-")
		g.secondaryKeyState = JimmyDoorDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyI:
		g.debugMessage = "Ignite Torch"
		g.addRowStr("Ignite Torch!")
		if !g.gameState.IgniteTorch() {
			g.addRowStr("None owned!")
		}
	case ebiten.KeyT:
		g.debugMessage = "Talk to..."
		g.addRowStr("Talk-")
		g.secondaryKeyState = TalkDirectionInput
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
		if !g.gameState.PartyState.Inventory.Provisions.Keys.HasSome() {
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
	case TalkDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			talk := g.smallMapTalkSecondary(getCurrentPressedArrowKeyAsDirection())
			if !talk {
				g.addRowStr("No-one to talk to!")
			}

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
	currentTile := g.gameState.MapState.LayeredMaps.GetTileRefByPosition(
		references2.SmallMapType,
		map_state.MapLayer,
		&g.gameState.MapState.PlayerLocation.Position,
		g.gameState.MapState.PlayerLocation.Floor)

	switch currentTile.Index {
	case indexes.AvatarOnLadderDown, indexes.LadderDown, indexes.Grate:
		if g.GetCurrentLocationReference().CanGoDownOneFloor(g.gameState.MapState.PlayerLocation.Floor) {
			g.gameState.MapState.PlayerLocation.Floor--
			g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
			g.output.AddRowStrWithTrim("Klimb-Down!")

			return
		} else {
			log.Fatal("Can't go lower my dude")
		}

	case indexes.AvatarOnLadderUp, indexes.LadderUp:
		if g.GetCurrentLocationReference().CanGoUpOneFloor(g.gameState.MapState.PlayerLocation.Floor) {
			g.gameState.MapState.PlayerLocation.Floor++
			g.gameState.UpdateSmallMap(g.gameReferences.TileReferences, g.gameReferences.LocationReferences)
			g.output.AddRowStrWithTrim("Klimb-Up!")

			return
		} else {
			log.Fatal("Can't go higher my dude")
		}
	}
	g.output.AddRowStrWithTrim("Klimb-")
	g.secondaryKeyState = KlimbDirectionInput
}

func (g *GameScene) smallMapKlimbSecondary(direction references2.Direction) {
	if !g.gameState.ActionKlimbSmallMap(direction) {
		g.output.AddRowStrWithTrim("What?")
	}
}

func (g *GameScene) smallMapPushSecondary(direction references2.Direction) {
	pushThingPos := direction.GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)

	pushThingTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(pushThingPos)

	if !pushThingTile.IsPushable {
		g.output.AddRowStrWithTrim("Won't budge!")
		return
	}

	if g.gameState.ActionPushSmallMap(direction) {
		// moved
	} else {
		// didn't move
	}
}

func (g *GameScene) smallMapOpenSecondary(direction references2.Direction) {
	openThingPos := direction.GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)
	openThingTile := g.gameState.GetLayeredMapByCurrentLocation().GetTileTopMapOnlyTile(openThingPos)

	if openThingTile.Index.IsDoor() {
		switch g.gameState.MapState.OpenDoor(direction) {
		case map_state.OpenDoorNotADoor:
			g.addRowStr("Bang to open!")
		case map_state.OpenDoorLocked:
			g.addRowStr("Locked!")
		case map_state.OpenDoorLockedMagical:
			g.addRowStr("Magically Locked!")
		case map_state.OpenDoorOpened:
			g.addRowStr("Opened!")
		default:
			log.Fatal("Unrecognized door open state")
		}
		return
	}

	openThingTopTile := g.gameState.GetLayeredMapByCurrentLocation().GetTopTile(openThingPos)
	if openThingTopTile.Index == indexes.Chest {
		if g.gameState.MapState.PlayerLocation.Location == references2.Lord_Britishs_Castle && g.gameState.MapState.PlayerLocation.Floor == references2.Basement {
			itemStack := references2.CreateNewItemStack(references2.LordBritishTreasure)
			g.output.AddRowStrWithTrim("Found:")
			g.output.AddRowStrWithTrim(g.gameReferences.InventoryItemReferences.GetListOfItems(&itemStack))
			g.gameState.ItemStacksMap.Push(openThingPos, &itemStack)
			g.gameState.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(*openThingPos)
			g.gameState.CurrentNPCAIController.FreshenExistingNPCsOnMap()
		}
	}
}

func (g *GameScene) smallMapJimmySecondary(direction references2.Direction) {
	jimmyResult := g.gameState.JimmyDoor(direction, &g.gameState.PartyState.Characters[0])

	switch jimmyResult {
	case gamestate.JimmyUnlocked:
		g.addRowStr("Unlocked!")
	case gamestate.JimmyNotADoor:
		g.addRowStr("Not lock!")
	case gamestate.JimmyBrokenPick, gamestate.JimmyLockedMagical:
		g.addRowStr("Key broke!")
	default:
		panic("unhandled default case")
	}
}

func (g *GameScene) smallMapGetSecondary(direction references2.Direction) {
	getThingPos := direction.GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)
	getThingTile := g.gameState.MapState.LayeredMaps.GetTileTopMapOnlyTileByPosition(references2.SmallMapType, getThingPos, g.gameState.MapState.PlayerLocation.Floor)
	mapLayers := g.gameState.MapState.LayeredMaps.GetLayeredMap(references2.SmallMapType, g.gameState.MapState.PlayerLocation.Floor)

	if g.gameState.ItemStacksMap.HasItemStackAtPosition(getThingPos) {
		item := g.gameState.ItemStacksMap.Pop(getThingPos)
		g.gameState.PartyState.Inventory.PutItemInInventory(item)

		itemRef := g.gameReferences.InventoryItemReferences.GetReferenceByItem(item.Item)
		// if item.
		g.addRowStr(fmt.Sprintf("%s!", itemRef.ItemName))
		return
	}

	switch getThingTile.Index {
	case indexes.WheatInField:
		g.addRowStr("Crops picked! Those aren't yours Avatar!")
		mapLayers.SetTileByLayer(map_state.MapLayer, getThingPos, indexes.PlowedField)
		g.gameState.PartyState.Karma.DecreaseKarma(1)
	case indexes.RightSconce, indexes.LeftScone:
		g.addRowStr("Borrowed!")
		g.gameState.PartyState.Inventory.Provisions.Torches.IncrementByOne()
		mapLayers.SetTileByLayer(map_state.MapLayer, getThingPos, indexes.BrickFloor)
	case indexes.TableFoodBoth, indexes.TableFoodBottom, indexes.TableFoodTop:
		if g.getFoodFromTable(direction) {
			g.addRowStr("Mmmmm...! But that food isn't yours!")
			g.gameState.PartyState.Inventory.Provisions.Food.IncrementByOne()
			g.gameState.PartyState.Karma.DecreaseKarma(1)
		}
	case indexes.Carpet2_MagicCarpet:
	}
}

func (g *GameScene) getFoodFromTable(direction references2.Direction) bool {
	getThingPos := direction.GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)
	getThingTile := g.gameState.MapState.LayeredMaps.GetTileTopMapOnlyTileByPosition(references2.SmallMapType, getThingPos, g.gameState.MapState.PlayerLocation.Floor)
	mapLayers := g.gameState.MapState.LayeredMaps.GetLayeredMap(references2.SmallMapType, g.gameState.MapState.PlayerLocation.Floor)

	var newTileIndex indexes.SpriteIndex

	switch direction {
	case references2.Down:
		if getThingTile.Index == indexes.TableFoodBoth {
			newTileIndex = indexes.TableFoodBottom
		} else if getThingTile.Index == indexes.TableFoodTop {
			newTileIndex = indexes.TableMiddle
		} else {
			return false
		}
	case references2.Up:
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

func (g *GameScene) smallMapTalkSecondary(direction references2.Direction) bool {
	talkThingPos := direction.GetNewPositionInDirection(&g.gameState.MapState.PlayerLocation.Position)
	npc := g.gameState.CurrentNPCAIController.GetNpcs().GetMapUnitAtPositionOrNil(*talkThingPos)

	if npc == nil {
		return false
	}

	if friendly, ok := (*npc).(*map_units.NPCFriendly); ok {
		talkDialog := NewTalkDialog(g, friendly.NPCReference)
		talkDialog.AddTestTest()
		g.dialogStack.PushModalDialog(talkDialog)
	}

	return true
}
