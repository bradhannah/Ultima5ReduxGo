package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameState) ActionTalkSmallMap(direction references.Direction) bool {
	talkThingPos := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	npc := g.CurrentNPCAIController.GetNpcs().GetMapUnitAtPositionOrNil(*talkThingPos)

	if npc == nil {
		g.SystemCallbacks.Message.AddRowStr("No-one to talk to!")
		return false
	}

	if friendly, ok := (*npc).(*map_units.NPCFriendly); ok {
		// TODO: Handle freed NPC acknowledgement (stocks/manacles with karma +2) here

		// Create and push dialog using dependency injection
		dialog := g.SystemCallbacks.Talk.CreateTalkDialog(friendly)
		if dialog != nil {
			g.SystemCallbacks.Talk.PushDialog(dialog)
			g.SystemCallbacks.Flow.AdvanceTime(1) // Talking takes time
			return true
		}
		return false
	}

	// NPC found but not friendly (hostile?)
	g.SystemCallbacks.Message.AddRowStr("Can't talk to that!")
	return false
}

func (g *GameState) ActionTalkLargeMap(direction references.Direction) bool {
	// Per Commands.md: Talk on overworld always returns "Funny, no response!"
	g.SystemCallbacks.Message.AddRowStr("Talk-Funny, no response!")
	g.SystemCallbacks.Flow.AdvanceTime(1)
	return true
}

func (g *GameState) ActionTalkCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Talk command - see Commands.md Talk section
	// Combat map variant of talk command - likely disabled during combat
	return true
}

func (g *GameState) ActionTalkDungeonMap(direction references.Direction) bool {
	// TODO: Implement dungeon map Talk command - see Commands.md Talk section
	// Dungeon map variant of talk command - may be limited
	return true
}
