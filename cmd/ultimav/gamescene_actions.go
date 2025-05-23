package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameScene) actionBoard() {
	vehicle := g.gameState.CurrentNPCAIController.GetNpcs().GetVehicleAtPositionOrNil(
		g.gameState.MapState.PlayerLocation.Position)

	if vehicle == nil {
		g.output.AddRowStr("Board what?")

		return
	}

	switch g.gameState.BoardVehicle(*vehicle) {
	case game_state.BoardVehicleResultSuccess:
		g.output.AddRowStr(vehicle.GetVehicleDetails().VehicleType.GetBoardString())
	case game_state.BoardVehicleResultNoVehicle:
		g.output.AddRowStr("Board what?")
	case game_state.BoardVehicleResultNoSkiffs:
		g.output.AddRowStr(vehicle.GetVehicleDetails().VehicleType.GetBoardString())
		g.output.AddRowStr("WARNING: NO SKIFFS ON BOARD!")
	}
}

func (g *GameScene) actionExit() {
	exittedVehicle := g.gameState.ExitVehicle()

	if exittedVehicle == nil || exittedVehicle.GetVehicleDetails().VehicleType == references.NoPartyVehicle {
		g.output.AddRowStr("X-it what?")

		return
	}

	g.output.AddRowStr(exittedVehicle.GetVehicleDetails().VehicleType.GetExitString())
}
