package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func (g *GameScene) smallMapInputHandler(key ebiten.Key) {
	switch key {
	case ebiten.KeyBackquote:
		g.bShowDebugConsole = !g.bShowDebugConsole
	case ebiten.KeyEnter:
		g.addRowStr("Enter")
	case ebiten.KeyUp:
		g.handleMovement(game_state.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(game_state.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(game_state.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(game_state.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyK:
		g.smallMapKlimb()
	case ebiten.KeyX:
		g.gameState.Location = references.Britannia_Underworld
		g.gameState.Floor = 0
		g.gameState.Position = g.gameState.LastLargeMapPosition
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
	}

	// only process end of turn if the turn is actually done.
	if g.secondaryKeyState == PrimaryInput {
		g.gameState.SmallMapProcessEndOfTurn()
	}
}

func (g *GameScene) smallMapHandleSecondaryInput() {
	arrowKey := getArrowKeyPressed()
	bIsArrowKeyPressed := arrowKey != nil

	switch g.secondaryKeyState {
	case JimmyDoorDirectionInput:
		if g.gameState.Provisions.QtyKeys <= 0 {
			g.addRowStr("No Keys!")
			g.secondaryKeyState = PrimaryInput
			g.keyboard.SetLastKeyPressedNow()
			return
		}

		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}
		g.appendDirectionToOutput()

		jimmyResult := g.gameState.JimmyDoor(getCurrentPressedArrowKeyAsDirection(), &g.gameState.Characters[0])

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

		g.secondaryKeyState = PrimaryInput

	case OpenDirectionInput:
		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}
		g.appendDirectionToOutput()

		switch g.gameState.OpenDoor(getCurrentPressedArrowKeyAsDirection()) {
		case game_state.OpenDoorNotADoor:
			g.addRowStr("Nothing to open!")
		case game_state.OpenDoorLocked:
			g.addRowStr("Locked!")
		case game_state.OpenDoorLockedMagical:
			g.addRowStr("Magically Locked!")
		case game_state.OpenDoorOpened:
			g.addRowStr("Opened!")
		default:
			log.Fatal("Unrecognized door open state")
		}

		g.secondaryKeyState = PrimaryInput
	case KlimbDirectionInput:
		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}
		g.appendDirectionToOutput()

		g.smallMapKlimbSecondary(getCurrentPressedArrowKeyAsDirection())

		g.secondaryKeyState = PrimaryInput
	case PushDirectionInput:
		if !bIsArrowKeyPressed {
			return
		}
		if !g.keyboard.TryToRegisterKeyPress(*arrowKey) {
			return
		}
		g.appendDirectionToOutput()

		g.smallMapPushSecondary(getCurrentPressedArrowKeyAsDirection())
		g.secondaryKeyState = PrimaryInput
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
	case indexes.AvatarOnLadderDown, indexes.Grate:
		if g.GetCurrentLocationReference().CanGoDownOneFloor(g.gameState.Floor) {
			g.gameState.Floor--
			g.output.AddRowStr("Klimb-Down!")
			return
		} else {
			log.Fatal("Can't go lower my dude")
		}

	case indexes.AvatarOnLadderUp:
		if g.GetCurrentLocationReference().CanGoUpOneFloor(g.gameState.Floor) {
			g.gameState.Floor++
			g.output.AddRowStr("Klimb-Up!")
			return
		} else {
			log.Fatal("Can't go higher my dude")
		}
	}
	g.output.AddRowStr("Klimb-")
	g.secondaryKeyState = KlimbDirectionInput
}

func (g *GameScene) smallMapKlimbSecondary(direction game_state.Direction) {
	if !g.gameState.KlimbSmallMap(direction) {
		g.output.AddRowStr("What?")
	}
}

func (g *GameScene) smallMapPushSecondary(direction game_state.Direction) {
	pushThingPos := direction.GetNewPositionInDirection(&g.gameState.Position)
	currentTile := g.gameState.LayeredMaps.GetTileRefByPosition(references.SmallMapType, game_state.MapLayer, pushThingPos, g.gameState.Floor)

	if !currentTile.IsPushable {
		g.output.AddRowStr("Won't budge!")
	}
	//switch spots
	g.gameState.LayeredMaps.GetLayeredMap(references.SmallMapType, g.gameState.Floor).SwapTiles(&g.gameState.Position, pushThingPos)
}
