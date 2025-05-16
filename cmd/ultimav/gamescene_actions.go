package main

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

func (g *GameScene) actionBoard() {
	vehicle := g.gameState.CurrentNPCAIController.GetNpcs().GetVehicleAtPositionOrNil(g.gameState.Position)
	if vehicle == nil {
		g.output.AddRowStr("Board what?")
		return
		//return game_state.NPCVehicle{VehicleType: references.NoPartyVehicle}
	}

	if !g.gameState.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(g.gameState.Position) {
		log.Fatalf("Unexpected - tried to remove NPC at position X=%d,Y=%d but failed", g.gameState.Position.X, g.gameState.Position.Y)
	}
	// have you already boarded something else?
	// I think it shouldn't matter because you shouldn't be allowed on a tile if you have boarded something
	// unless you are a skiff on a frigate

	if !g.gameState.BoardVehicle(*vehicle) {
		return
		//return game_state.NPCVehicle{VehicleType: references.NoPartyVehicle}
	}

	g.output.AddRowStr(vehicle.VehicleType.GetBoardString())

	// we are boarded...
	return
	//return vehicle
}

// func (g *GameScene) actionBoard() game_state.NPCVehicle {
// 	getThingTile := g.gameState.GetLayeredMapByCurrentLocation().GetTileByLayer(game_state.MapUnitLayer, &g.gameState.Position)

// 	vehicle := references.GetVehicleFromSpriteIndex(getThingTile.Index)

// 	if vehicle == references.NoPartyVehicle {
// 		g.output.AddRowStr("Board what?")
// 		return game_state.NPCVehicle{VehicleType: references.NoPartyVehicle}
// 	}

// 	if !g.gameState.BoardVehicle(vehicle) {
// 		return game_state.NPCVehicle{VehicleType: references.NoPartyVehicle}
// 	}

// 	if !g.gameState.CurrentNPCAIController.GetNpcs().RemoveNPCAtPosition(g.gameState.Position) {
// 		log.Fatalf("Unexpected - tried to remove NPC at position X=%d,Y=%d but failed", g.gameState.Position.X, g.gameState.Position.Y)
// 	}

// 	g.output.AddRowStr(vehicle.GetBoardString())

// 	return vehicle
// }

func (g *GameScene) actionExit() {
	vr := g.gameState.ExitVehicle()

	if vr.ExittedVehicle.VehicleType == references.NoPartyVehicle {
		g.output.AddRowStr("X-it what?")
		return
	}
	g.output.AddRowStr("X-it-")
	switch vr.ExittedVehicle.VehicleType {
	case references.FrigateVehicle:
		g.output.AppendToCurrentRowStr("frigate!")
	case references.CarpetVehicle:
		g.output.AppendToCurrentRowStr("carpet!")
	case references.SkiffVehicle:
		g.output.AppendToCurrentRowStr("skiff!")
	case references.HorseVehicle:
		g.output.AddRowStr("horse!")
	}
}
