package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

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
	case ebiten.KeyS:
		g.debugMessage = "Search"
		g.addRowStr("Search-")
		g.secondaryKeyState = SearchDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyA:
		g.debugMessage = "Attack"
		g.addRowStr("Attack-")
		g.secondaryKeyState = AttackDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyU:
		g.debugMessage = "Use"
		g.addRowStr("Use-")
		g.secondaryKeyState = UseDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyY:
		g.debugMessage = "Yell"
		g.addRowStr("Yell-")
		g.secondaryKeyState = YellDirectionInput
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
	case SearchDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapSearchSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case AttackDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapAttackSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case UseDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapUseSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case YellDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.smallMapYellSecondary(getCurrentPressedArrowKeyAsDirection())
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

	// Early validation - avoid GameState call if obviously invalid
	if !g.gameState.IsPushable(pushThingTile) {
		g.output.AddRowStrWithTrim("Won't budge!") // Direct UI - no game logic
		return
	}

	// Delegate everything else to GameState - it handles all feedback via SystemCallbacks
	g.gameState.ActionPushSmallMap(direction)
}

func (g *GameScene) smallMapOpenSecondary(direction references2.Direction) {
	// Delegate all logic to GameState - it handles all feedback via SystemCallbacks
	g.gameState.ActionOpenSmallMap(direction)
}

func (g *GameScene) smallMapJimmySecondary(direction references2.Direction) {
	// Delegate all logic to GameState - it handles all feedback via SystemCallbacks
	g.gameState.ActionJimmySmallMap(direction)
}

func (g *GameScene) smallMapGetSecondary(direction references2.Direction) {
	// Delegate all logic to GameState - it handles all feedback via SystemCallbacks
	g.gameState.ActionGetSmallMap(direction)
}

func (g *GameScene) smallMapTalkSecondary(direction references2.Direction) bool {
	// Get detailed result from GameState
	talkResult := g.gameState.ActionTalkSmallMap(direction)

	if !talkResult.Success {
		if talkResult.Message != "" {
			g.addRowStr(talkResult.Message)
		}
		return false
	}

	if talkResult.NPC != nil {
		// Use linear conversation system - UI responsibility
		linearTalkDialog := NewLinearTalkDialog(g, talkResult.NPC.NPCReference)
		linearTalkDialog.AddTestTest()
		g.dialogStack.PushModalDialog(linearTalkDialog)
		return true
	}

	return false
}

func (g *GameScene) smallMapSearchSecondary(direction references2.Direction) {
	// TODO: Implement Search secondary action
	success := g.gameState.ActionSearchSmallMap(direction)
	if !success {
		g.addRowStr("Nothing found!")
	}
}

func (g *GameScene) smallMapAttackSecondary(direction references2.Direction) {
	// TODO: Implement Attack secondary action
	success := g.gameState.ActionAttackSmallMap(direction)
	if !success {
		g.addRowStr("Nothing to attack!")
	}
}

func (g *GameScene) smallMapUseSecondary(direction references2.Direction) {
	// TODO: Implement Use secondary action
	success := g.gameState.ActionUseSmallMap(direction)
	if !success {
		g.addRowStr("Cannot use!")
	}
}

func (g *GameScene) smallMapYellSecondary(direction references2.Direction) {
	// TODO: Implement Yell secondary action
	success := g.gameState.ActionYellSmallMap(direction)
	if !success {
		g.addRowStr("No effect!")
	}
}
