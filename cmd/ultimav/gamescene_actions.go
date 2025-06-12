package main

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

func (g *GameScene) actionBoard() {
	vehicle := g.gameState.CurrentNPCAIController.GetNpcs().GetVehicleAtPositionOrNil(
		g.gameState.MapState.PlayerLocation.Position)

	if vehicle == nil {
		g.output.AddRowStrWithTrim("Board what?")

		return
	}

	switch g.gameState.BoardVehicle(*vehicle) {
	case game_state.BoardVehicleResultSuccess:
		g.output.AddRowStrWithTrim(vehicle.GetVehicleDetails().VehicleType.GetBoardString())
	case game_state.BoardVehicleResultNoVehicle:
		g.output.AddRowStrWithTrim("Board what?")
	case game_state.BoardVehicleResultNoSkiffs:
		g.output.AddRowStrWithTrim(vehicle.GetVehicleDetails().VehicleType.GetBoardString())
		g.output.AddRowStrWithTrim("WARNING: NO SKIFFS ON BOARD!")
	}
}

func (g *GameScene) actionExit() {
	exittedVehicle := g.gameState.ExitVehicle()

	if exittedVehicle == nil || exittedVehicle.GetVehicleDetails().VehicleType == references.NoPartyVehicle {
		g.output.AddRowStrWithTrim("X-it what?")

		return
	}

	g.output.AddRowStrWithTrim(exittedVehicle.GetVehicleDetails().VehicleType.GetExitString())
}
