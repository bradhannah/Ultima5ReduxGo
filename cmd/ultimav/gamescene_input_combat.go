package main

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *GameScene) combatMapInputHandler(key ebiten.Key) {
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			// TODO: Implement combat quick exit when combat system is available
			return
		}
	}

	switch key {
	case ebiten.KeyEscape:
		g.DoEscapeMenu()
		return
	case ebiten.KeySpace:
		g.addRowStr("Pass")
		// TODO: Implement combat turn passing
		// g.gameState.PassCombatTurn()
	case ebiten.KeyBackquote:
		g.toggleDebug()
		return
	case ebiten.KeyUp:
		g.handleMovement(references2.Up.GetDirectionCompassName(), ebiten.KeyUp)
	case ebiten.KeyDown:
		g.handleMovement(references2.Down.GetDirectionCompassName(), ebiten.KeyDown)
	case ebiten.KeyLeft:
		g.handleMovement(references2.Left.GetDirectionCompassName(), ebiten.KeyLeft)
	case ebiten.KeyRight:
		g.handleMovement(references2.Right.GetDirectionCompassName(), ebiten.KeyRight)
	case ebiten.KeyL:
		g.addRowStr("Look-")
		g.secondaryKeyState = LookDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyG:
		g.addRowStr("Get-")
		g.secondaryKeyState = GetDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyP:
		g.addRowStr("Push-")
		g.secondaryKeyState = PushDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyO:
		g.addRowStr("Open-")
		g.secondaryKeyState = OpenDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyJ:
		g.addRowStr("Jimmy-")
		g.secondaryKeyState = JimmyDoorDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyK:
		g.addRowStr("Klimb-")
		g.secondaryKeyState = KlimbDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyI:
		g.addRowStr("Ignite Torch!")
		if !g.gameState.IgniteTorch() {
			g.addRowStr("None owned!")
		}
	case ebiten.KeyT:
		g.addRowStr("Talk-")
		g.secondaryKeyState = TalkDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyS:
		g.addRowStr("Search-")
		g.secondaryKeyState = SearchDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyA:
		g.addRowStr("Attack-")
		g.secondaryKeyState = AttackDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyU:
		g.addRowStr("Use-")
		g.secondaryKeyState = UseDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyY:
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

func (g *GameScene) combatMapHandleSecondaryInput() {
	switch g.secondaryKeyState {
	case JimmyDoorDirectionInput:
		if !g.gameState.PartyState.Inventory.Provisions.Keys.HasSome() {
			g.addRowStr("No Keys!")
			g.secondaryKeyState = PrimaryInput
			g.keyboard.SetLastKeyPressedNow()
			return
		}

		if g.isDirectionKeyValidAndOutput() {
			g.combatMapJimmySecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case OpenDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapOpenSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case KlimbDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapKlimbSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case PushDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapPushSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case GetDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapGetSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case LookDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapLookSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case TalkDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			success := g.combatMapTalkSecondary(getCurrentPressedArrowKeyAsDirection())
			if !success {
				g.addRowStr("Not in combat!")
			}
			g.secondaryKeyState = PrimaryInput
		}
	case SearchDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapSearchSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case AttackDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapAttackSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case UseDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapUseSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case YellDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.combatMapYellSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	}

	// only process end of turn if the turn is actually done.
	if g.secondaryKeyState == PrimaryInput {
		g.gameState.FinishTurn()
	}
}

// Combat Map Secondary Action Handlers

func (g *GameScene) combatMapKlimbSecondary(direction references2.Direction) {
	success := g.gameState.ActionKlimbCombatMap(direction)
	if !success {
		g.addRowStr("Nothing to klimb!")
	}
}

func (g *GameScene) combatMapPushSecondary(direction references2.Direction) {
	success := g.gameState.ActionPushCombatMap(direction)
	if !success {
		g.addRowStr("Won't budge!")
	}
}

func (g *GameScene) combatMapOpenSecondary(direction references2.Direction) {
	success := g.gameState.ActionOpenCombatMap(direction)
	if !success {
		g.addRowStr("Nothing to open!")
	}
}

func (g *GameScene) combatMapJimmySecondary(direction references2.Direction) {
	success := g.gameState.ActionJimmyCombatMap(direction)
	if !success {
		g.addRowStr("Nothing to jimmy!")
	}
}

func (g *GameScene) combatMapGetSecondary(direction references2.Direction) {
	success := g.gameState.ActionGetCombatMap(direction)
	if !success {
		g.addRowStr("Nothing to get!")
	}
}

func (g *GameScene) combatMapLookSecondary(direction references2.Direction) {
	success := g.gameState.ActionLookCombatMap(direction)
	if !success {
		g.addRowStr("Nothing to see!")
	}
}

func (g *GameScene) combatMapTalkSecondary(direction references2.Direction) bool {
	success := g.gameState.ActionTalkCombatMap(direction)
	return success
}

func (g *GameScene) combatMapSearchSecondary(direction references2.Direction) {
	success := g.gameState.ActionSearchCombatMap(direction)
	if !success {
		g.addRowStr("Nothing found!")
	}
}

func (g *GameScene) combatMapAttackSecondary(direction references2.Direction) {
	success := g.gameState.ActionAttackCombatMap(direction)
	if !success {
		g.addRowStr("Nothing to attack!")
	}
}

func (g *GameScene) combatMapUseSecondary(direction references2.Direction) {
	success := g.gameState.ActionUseCombatMap(direction)
	if !success {
		g.addRowStr("Can't use that!")
	}
}

func (g *GameScene) combatMapYellSecondary(direction references2.Direction) {
	success := g.gameState.ActionYellCombatMap(direction)
	if !success {
		g.addRowStr("No effect!")
	}
}
