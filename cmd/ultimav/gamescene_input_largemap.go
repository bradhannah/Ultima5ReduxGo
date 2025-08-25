package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameScene) largeMapInputHandler(key ebiten.Key) {
	switch key {
	case ebiten.KeyEscape:
		g.DoEscapeMenu()

		return
	case ebiten.KeySpace:
		g.gameState.ActionPass()
	case ebiten.KeyBackquote:
		g.toggleDebug()
		return
	case ebiten.KeyEnter:
		g.gameState.ActionEnterInput()
	case ebiten.KeyUp:
		g.handleMovement(references2.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(references2.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(references2.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(references2.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyB:
		g.gameState.ActionBoard()
	case ebiten.KeyK:
		g.addRowStr("Klimb-")
		g.secondaryKeyState = KlimbDirectionInput
	case ebiten.KeyG:
		g.addRowStr("Get-")
		g.secondaryKeyState = GetDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyX:
		g.gameState.ActionExit()
	case ebiten.KeyP:
		g.addRowStr("Push-")
		g.secondaryKeyState = PushDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyL:
		g.addRowStr("Look-")
		g.secondaryKeyState = LookDirectionInput
	case ebiten.KeyE:
		g.debugMessage = "Enter a place"
		g.gameState.ActionEnterLargeMap()

	case ebiten.KeyO:
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
		g.gameState.ActionIgnite()
	case ebiten.KeyT:
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

func (g *GameScene) largeMapHandleSecondaryInput() {
	switch g.secondaryKeyState {
	case KlimbDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.gameState.ActionKlimbLargeMap(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case LookDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.gameState.ActionLookLargeMap(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case JimmyDoorDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.largeMapJimmySecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case GetDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.gameState.ActionGetLargeMap(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case PushDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.gameState.ActionPushLargeMap(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case OpenDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.gameState.ActionOpenLargeMap(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case TalkDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.gameState.ActionTalkLargeMap(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	default:
		// better safe than sorry
		g.secondaryKeyState = PrimaryInput
	}
}

func (g *GameScene) largeMapJimmySecondary(direction references2.Direction) {
	// Delegate all logic to GameState - it handles all feedback via SystemCallbacks
	g.gameState.ActionJimmyLargeMap(direction)
}
