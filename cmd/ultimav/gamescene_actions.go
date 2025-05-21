package main

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameScene) actionBoard() {
	vehicle := g.gameState.CurrentNPCAIController.GetNpcs().GetVehicleAtPositionOrNil(g.gameState.MapState.PlayerLocation.Position)
	if vehicle == nil {
		g.output.AddRowStr("Board what?")
		return
	}

	// if frigate under, in skiff, add to skiff, board ship, delete skiff

	if !g.gameState.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(g.gameState.MapState.PlayerLocation.Position) {
		log.Fatalf("Unexpected - tried to remove NPC at position X=%d,Y=%d but failed", g.gameState.MapState.PlayerLocation.Position.X, g.gameState.MapState.PlayerLocation.Position.Y)
	}
	// have you already boarded something else?
	// I think it shouldn't matter because you shouldn't be allowed on a tile if you have boarded something
	// unless you are a skiff on a frigate

	if !g.gameState.BoardVehicle(*vehicle) {
		return
	}

	g.output.AddRowStr(vehicle.GetVehicleDetails().VehicleType.GetBoardString())

	// we are boarded...
	return
	// return vehicle
}

func (g *GameScene) actionExit() {
	exittedVehicle := g.gameState.ExitVehicle()

	if exittedVehicle == nil || exittedVehicle.GetVehicleDetails().VehicleType == references.NoPartyVehicle {
		g.output.AddRowStr("X-it what?")
		return
	}

	g.output.AddRowStr(exittedVehicle.GetVehicleDetails().VehicleType.GetExitString())
}
