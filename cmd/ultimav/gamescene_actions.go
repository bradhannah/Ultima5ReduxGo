package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameScene) actionBoard() references.PartyVehicle {
	getThingTile := g.gameState.GetLayeredMapByCurrentLocation().GetTileByLayer(game_state.MapUnitLayer, &g.gameState.Position)

	vehicle := references.GetVehicleFromSpriteIndex(getThingTile.Index)

	if vehicle == references.NoPartyVehicle {
		g.output.AddRowStr("Board what?")
		return references.NoPartyVehicle
	}

	if !g.gameState.BoardVehicle(vehicle) {
		return references.NoPartyVehicle
	}

	g.gameState.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(g.gameState.Position)

	switch vehicle {
	case references.CarpetVehicle:
		g.output.AddRowStr("Board carpet")
	case references.HorseVehicle:
		g.output.AddRowStr("Board horse")
	case references.FrigateVehicle:
		g.output.AddRowStr("Board frigate")
	case references.SkiffVehicle:
		g.output.AddRowStr("Board skiff")
	default:
	}

	return vehicle
}

func (g *GameScene) actionExit() {
	vr := g.gameState.ExitVehicle()

	if vr.ExittedVehicle == references.NoPartyVehicle {
		g.output.AddRowStr("X-it what?")
		return
	}
	switch vr.ExittedVehicle {
	case references.CarpetVehicle:
	case references.FrigateVehicle:
	case references.SkiffVehicle:
	case references.HorseVehicle:
		g.output.AddRowStr("Boop?")
	}
}
