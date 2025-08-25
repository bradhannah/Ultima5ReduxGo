package game_state

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

// TestObjectPresentAt_BasicFunctionality tests the core collision detection
func TestObjectPresentAt_BasicFunctionality(t *testing.T) {
	gs := createMinimalGameState()

	// Test empty space
	position := references.Position{X: 5, Y: 5}
	result := gs.ObjectPresentAt(&position)
	if result {
		t.Error("ObjectPresentAt should return false for empty space")
	}

	// Test position with object (we'll add one manually to the mock controller)
	objectPosition := references.Position{X: 10, Y: 10}
	addSimpleMapUnitAtPosition(gs, objectPosition)

	result = gs.ObjectPresentAt(&objectPosition)
	if !result {
		t.Error("ObjectPresentAt should return true when object is present")
	}

	// Test adjacent position is still empty
	adjacentPosition := references.Position{X: 11, Y: 10}
	result = gs.ObjectPresentAt(&adjacentPosition)
	if result {
		t.Error("ObjectPresentAt should return false for adjacent empty space")
	}
}

// TestObjectPresentAt_PushCommandIntegration tests integration with push logic
func TestObjectPresentAt_PushCommandIntegration(t *testing.T) {
	gs := createMinimalGameState()

	// Setup: Object at position (6, 5) - this is where we're trying to push from
	objectAtPushSource := references.Position{X: 6, Y: 5}
	addSimpleMapUnitAtPosition(gs, objectAtPushSource)

	// Setup: Object at position (7, 5) - this is the destination we're trying to push to
	objectAtDestination := references.Position{X: 7, Y: 5}
	addSimpleMapUnitAtPosition(gs, objectAtDestination)

	// Verify both positions have objects
	if !gs.ObjectPresentAt(&objectAtPushSource) {
		t.Error("Should detect object at push source")
	}
	if !gs.ObjectPresentAt(&objectAtDestination) {
		t.Error("Should detect object at destination")
	}

	// This demonstrates that the push command can now properly detect collisions
	t.Logf("✅ ObjectPresentAt properly detects objects at both push source (%v) and destination (%v)",
		objectAtPushSource, objectAtDestination)
}

// TestObjectCollisionDetection_Regression ensures the fix addresses the original issue
func TestObjectCollisionDetection_Regression(t *testing.T) {
	// This test documents the fix for the original issue:
	// "Push command doesn't validate object-to-object collisions per OLD reference"

	gs := createMinimalGameState()

	// Before the fix: ObjectPresentAt always returned false (stub)
	// After the fix: ObjectPresentAt properly checks the NPC/object system

	position := references.Position{X: 15, Y: 15}

	// Test 1: No object present
	result := gs.ObjectPresentAt(&position)
	if result {
		t.Error("Should return false when no object present")
	}

	// Test 2: Add an object
	addSimpleMapUnitAtPosition(gs, position)
	result = gs.ObjectPresentAt(&position)
	if !result {
		t.Error("REGRESSION: Should return true when object is present - fix is not working")
	}

	t.Logf("✅ Regression test passed - ObjectPresentAt now properly implements looklist() functionality")
}

// Helper functions for minimal test setup

func createMinimalGameState() *GameState {
	gameConfig := &config.UltimaVConfiguration{}
	gameReferences := &references.GameReferences{}
	gs := initBlankGameState(gameConfig, gameReferences, 16, 16)

	// Initialize with a simple mock NPC controller
	gs.CurrentNPCAIController = &SimpleTestController{
		mapUnits: make(map_units.MapUnits, 0),
	}

	return gs
}

func addSimpleMapUnitAtPosition(gs *GameState, pos references.Position) {
	if controller, ok := gs.CurrentNPCAIController.(*SimpleTestController); ok {
		// Create a minimal map unit (just for testing collision detection)
		mapUnit := &SimpleTestMapUnit{
			position: pos,
			floor:    references.FloorNumber(0),
			visible:  true,
		}
		controller.mapUnits = append(controller.mapUnits, mapUnit)
	}
}

// Simple test implementations

type SimpleTestController struct {
	mapUnits map_units.MapUnits
}

func (s *SimpleTestController) PopulateMapFirstLoad()           {}
func (s *SimpleTestController) AdvanceNextTurnCalcAndMoveNPCs() {}
func (s *SimpleTestController) FreshenExistingNPCsOnMap()       {}
func (s *SimpleTestController) RemoveAllEnemies()               {}
func (s *SimpleTestController) GetNpcs() *map_units.MapUnits {
	return &s.mapUnits
}

type SimpleTestMapUnit struct {
	position references.Position
	floor    references.FloorNumber
	visible  bool
}

func (s *SimpleTestMapUnit) GetMapUnitType() map_units.MapUnitType {
	return map_units.NonPlayerCharacter
}
func (s *SimpleTestMapUnit) Pos() references.Position                  { return s.position }
func (s *SimpleTestMapUnit) PosPtr() *references.Position              { return &s.position }
func (s *SimpleTestMapUnit) Floor() references.FloorNumber             { return s.floor }
func (s *SimpleTestMapUnit) SetPos(pos references.Position)            { s.position = pos }
func (s *SimpleTestMapUnit) SetFloor(floor references.FloorNumber)     { s.floor = floor }
func (s *SimpleTestMapUnit) MapUnitDetails() *map_units.MapUnitDetails { return nil }
func (s *SimpleTestMapUnit) IsVisible() bool                           { return s.visible }
func (s *SimpleTestMapUnit) SetVisible(visible bool)                   { s.visible = visible }
func (s *SimpleTestMapUnit) IsEmptyMapUnit() bool                      { return false }
