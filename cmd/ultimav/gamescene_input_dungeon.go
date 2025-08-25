package main

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *GameScene) dungeonMapInputHandler(key ebiten.Key) {
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		if ebiten.IsKeyPressed(ebiten.KeyX) {
			// TODO: Implement dungeon quick exit when dungeon system is available
			return
		}
	}

	switch key {
	case ebiten.KeyEscape:
		g.DoEscapeMenu()
		return
	case ebiten.KeySpace:
		g.addRowStr("Pass")
		// TODO: Implement dungeon turn passing with hazard/lighting updates
		// g.gameState.PassDungeonTurn()
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
		g.gameState.ActionIgnite()
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
	case ebiten.KeyF:
		g.addRowStr("Fire-")
		g.secondaryKeyState = FireDirectionInput
		g.keyboard.SetAllowKeyPressImmediately()
	case ebiten.KeyR:
		g.addRowStr("Ready...")
		g.gameState.ActionReadyDungeonMap()
	case ebiten.KeyV:
		g.addRowStr("View...")
		g.gameState.ActionViewDungeonMap()
	case ebiten.KeyZ:
		g.addRowStr("Ztats...")
		g.gameState.ActionZtatsDungeonMap()
	case ebiten.KeyM:
		g.addRowStr("Mix...")
		g.gameState.ActionMixDungeonMap()
	case ebiten.KeyC:
		g.addRowStr("Cast...")
		g.gameState.ActionCastDungeonMap()
	case ebiten.KeyN:
		g.addRowStr("New Order...")
		g.gameState.ActionNewOrderDungeonMap()
	case ebiten.KeyH:
		g.addRowStr("Hole up & camp...")
		g.gameState.ActionHoleUpDungeonMap()
	case ebiten.KeyQ:
		g.addRowStr("Escape...")
		g.gameState.ActionEscapeDungeonMap()
	default:
		return
	}

	// only process end of turn if the turn is actually done.
	if g.secondaryKeyState == PrimaryInput {
		g.gameState.FinishTurn()
	}
}

func (g *GameScene) dungeonMapHandleSecondaryInput() {
	switch g.secondaryKeyState {
	case JimmyDoorDirectionInput:
		if !g.gameState.PartyState.Inventory.Provisions.Keys.HasSome() {
			g.addRowStr("No Keys!")
			g.secondaryKeyState = PrimaryInput
			g.keyboard.SetLastKeyPressedNow()
			return
		}

		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapJimmySecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case OpenDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapOpenSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case KlimbDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapKlimbSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case PushDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapPushSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case GetDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapGetSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case LookDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapLookSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case TalkDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			success := g.dungeonMapTalkSecondary(getCurrentPressedArrowKeyAsDirection())
			if !success {
				g.addRowStr("No-one to talk to!")
			}
			g.secondaryKeyState = PrimaryInput
		}
	case SearchDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapSearchSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case AttackDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapAttackSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case UseDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapUseSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case YellDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapYellSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	case FireDirectionInput:
		if g.isDirectionKeyValidAndOutput() {
			g.dungeonMapFireSecondary(getCurrentPressedArrowKeyAsDirection())
			g.secondaryKeyState = PrimaryInput
		}
	}

	// only process end of turn if the turn is actually done.
	if g.secondaryKeyState == PrimaryInput {
		g.gameState.FinishTurn()
	}
}

// Dungeon Map Secondary Action Handlers

func (g *GameScene) dungeonMapKlimbSecondary(direction references2.Direction) {
	success := g.gameState.ActionKlimbDungeonMap(direction)
	if !success {
		g.addRowStr("Nothing to klimb!")
	}
}

func (g *GameScene) dungeonMapPushSecondary(direction references2.Direction) {
	success := g.gameState.ActionPushDungeonMap(direction)
	if !success {
		g.addRowStr("Won't budge!")
	}
}

func (g *GameScene) dungeonMapOpenSecondary(direction references2.Direction) {
	success := g.gameState.ActionOpenDungeonMap(direction)
	if !success {
		g.addRowStr("Nothing to open!")
	}
}

func (g *GameScene) dungeonMapJimmySecondary(direction references2.Direction) {
	success := g.gameState.ActionJimmyDungeonMap(direction)
	if !success {
		g.addRowStr("Nothing to jimmy!")
	}
}

func (g *GameScene) dungeonMapGetSecondary(direction references2.Direction) {
	success := g.gameState.ActionGetDungeonMap(direction)
	if !success {
		g.addRowStr("Nothing to get!")
	}
}

func (g *GameScene) dungeonMapLookSecondary(direction references2.Direction) {
	success := g.gameState.ActionLookDungeonMap(direction)
	if !success {
		g.addRowStr("Too dark to see!")
	}
}

func (g *GameScene) dungeonMapTalkSecondary(direction references2.Direction) bool {
	success := g.gameState.ActionTalkDungeonMap(direction)
	return success
}

func (g *GameScene) dungeonMapSearchSecondary(direction references2.Direction) {
	success := g.gameState.ActionSearchDungeonMap(direction)
	if !success {
		g.addRowStr("Nothing found!")
	}
}

func (g *GameScene) dungeonMapAttackSecondary(direction references2.Direction) {
	success := g.gameState.ActionAttackDungeonMap(direction)
	if !success {
		g.addRowStr("Nothing to attack!")
	}
}

func (g *GameScene) dungeonMapUseSecondary(direction references2.Direction) {
	success := g.gameState.ActionUseDungeonMap(direction)
	if !success {
		g.addRowStr("Can't use that!")
	}
}

func (g *GameScene) dungeonMapYellSecondary(direction references2.Direction) {
	success := g.gameState.ActionYellDungeonMap(direction)
	if !success {
		g.addRowStr("No effect!")
	}
}

func (g *GameScene) dungeonMapFireSecondary(direction references2.Direction) {
	success := g.gameState.ActionFireDungeonMap(direction)
	if !success {
		g.addRowStr("Can't fire here!")
	}
}
