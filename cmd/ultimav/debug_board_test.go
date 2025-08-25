package main

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TestDebugBoardVehicle_Horse tests that the debug board function actually boards a horse
func TestDebugBoardVehicle_Horse(t *testing.T) {
	// Create a minimal GameScene for testing
	gameScene := &GameScene{}

	// Initialize minimal game state
	gameScene.gameState = &game_state.GameState{}

	// Initialize PartyVehicle to "no vehicle" state
	gameScene.gameState.PartyVehicle = map_units.NewNPCFriendlyVehiceNoVehicle()

	// Set a basic player position
	gameScene.gameState.MapState.PlayerLocation.Position = references.Position{X: 10, Y: 10}
	gameScene.gameState.MapState.PlayerLocation.Floor = 0

	// Verify player starts on foot
	initialVehicleType := gameScene.gameState.PartyVehicle.GetVehicleDetails().VehicleType
	if initialVehicleType != references.NoPartyVehicle {
		t.Errorf("Expected player to start on foot, got vehicle type: %v", initialVehicleType)
	}

	// Call debug board function for horse
	result := gameScene.debugBoardVehicle(references.HorseVehicle)

	// Verify the function returned success
	if !result {
		t.Error("debugBoardVehicle should have returned true for horse")
	}

	// Verify NPCType is correctly set to Vehicle
	actualNPCType := gameScene.gameState.PartyVehicle.NPCReference.GetNPCType()
	if actualNPCType != references.Vehicle {
		t.Errorf("Expected NPCType to be Vehicle (%v), got %v", references.Vehicle, actualNPCType)
	}

	// Verify player is now on a horse
	finalVehicleType := gameScene.gameState.PartyVehicle.GetVehicleDetails().VehicleType
	if finalVehicleType != references.HorseVehicle {
		t.Errorf("Expected player to be on horse, got vehicle type: %v", finalVehicleType)
	}

	// Verify the boarded sprite is valid
	boardedSprite := gameScene.gameState.PartyVehicle.GetVehicleDetails().GetBoardedSpriteIndex()
	if boardedSprite == 0 {
		t.Error("Boarded sprite index should not be zero/empty")
	}
}

// TestDebugBoardVehicle_AllTypes tests all vehicle types
func TestDebugBoardVehicle_AllTypes(t *testing.T) {
	vehicleTypes := []references.VehicleType{
		references.HorseVehicle,
		references.CarpetVehicle,
		references.SkiffVehicle,
		references.FrigateVehicle,
	}

	for _, vehicleType := range vehicleTypes {
		t.Run(vehicleType.GetLowerCaseName(), func(t *testing.T) {
			// Create a minimal GameScene for testing
			gameScene := &GameScene{}

			// Initialize minimal game state
			gameScene.gameState = &game_state.GameState{}

			// Initialize PartyVehicle to "no vehicle" state
			gameScene.gameState.PartyVehicle = map_units.NewNPCFriendlyVehiceNoVehicle()

			// Set a basic player position
			gameScene.gameState.MapState.PlayerLocation.Position = references.Position{X: 10, Y: 10}
			gameScene.gameState.MapState.PlayerLocation.Floor = 0

			// Verify player starts on foot
			initialVehicleType := gameScene.gameState.PartyVehicle.GetVehicleDetails().VehicleType
			if initialVehicleType != references.NoPartyVehicle {
				t.Errorf("Expected player to start on foot, got vehicle type: %v", initialVehicleType)
			}

			// Call debug board function
			result := gameScene.debugBoardVehicle(vehicleType)

			// Verify the function returned success
			if !result {
				t.Errorf("debugBoardVehicle should have returned true for %s", vehicleType.GetLowerCaseName())
			}

			// Verify NPCType is correctly set to Vehicle
			actualNPCType := gameScene.gameState.PartyVehicle.NPCReference.GetNPCType()
			if actualNPCType != references.Vehicle {
				t.Errorf("Expected NPCType to be Vehicle (%v), got %v for %s", references.Vehicle, actualNPCType, vehicleType.GetLowerCaseName())
			}

			// Verify player is now on the correct vehicle
			finalVehicleType := gameScene.gameState.PartyVehicle.GetVehicleDetails().VehicleType
			if finalVehicleType != vehicleType {
				t.Errorf("Expected player to be on %s, got vehicle type: %v", vehicleType.GetLowerCaseName(), finalVehicleType)
			}

			// Verify the boarded sprite is valid
			boardedSprite := gameScene.gameState.PartyVehicle.GetVehicleDetails().GetBoardedSpriteIndex()
			if boardedSprite == 0 {
				t.Errorf("%s boarded sprite index should not be zero/empty", vehicleType.GetLowerCaseName())
			}

			// Special verification for frigates - should have skiffs
			if vehicleType == references.FrigateVehicle {
				if !gameScene.gameState.PartyVehicle.GetVehicleDetails().HasAtLeastOneSkiff() {
					t.Error("Frigate should have at least one skiff")
				}
			}
		})
	}
}
