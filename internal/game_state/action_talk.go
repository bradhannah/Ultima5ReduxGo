package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TalkResult represents the result of attempting to talk
type TalkResult struct {
	Success bool
	NPC     *map_units.NPCFriendly // nil if no NPC found or not friendly
	Message string                 // Error message if any
}

func (g *GameState) ActionTalkSmallMap(direction references.Direction) TalkResult {
	talkThingPos := direction.GetNewPositionInDirection(&g.MapState.PlayerLocation.Position)
	npc := g.CurrentNPCAIController.GetNpcs().GetMapUnitAtPositionOrNil(*talkThingPos)

	if npc == nil {
		return TalkResult{
			Success: false,
			NPC:     nil,
			Message: "No-one to talk to!",
		}
	}

	if friendly, ok := (*npc).(*map_units.NPCFriendly); ok {
		// TODO: Handle freed NPC acknowledgement (stocks/manacles with karma +2) here
		// For now, return the friendly NPC for the UI to handle dialog
		g.SystemCallbacks.Flow.AdvanceTime(1) // Talking takes time
		return TalkResult{
			Success: true,
			NPC:     friendly,
			Message: "",
		}
	}

	// NPC found but not friendly (hostile?)
	return TalkResult{
		Success: false,
		NPC:     nil,
		Message: "Can't talk to that!",
	}
}

func (g *GameState) ActionTalkLargeMap(direction references.Direction) bool {
	// TODO: Implement large map Talk command - see Commands.md Talk section
	// Large map variant of talk command
	return true
}

func (g *GameState) ActionTalkCombatMap(direction references.Direction) bool {
	// TODO: Implement combat map Talk command - see Commands.md Talk section
	// Combat map variant of talk command - likely disabled during combat
	return true
}
