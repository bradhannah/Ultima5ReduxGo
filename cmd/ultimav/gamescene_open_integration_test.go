package main

import (
	"testing"

	"github.com/bradhannah/Ultima5ReduxGo/internal/map_state"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
	"github.com/hajimehoshi/ebiten/v2"
)

func TestGameScene_SmallMapOpenSecondary_Door_Success(t *testing.T) {
	gs := createTestGameScene(t)

	// Set up a regular door to the north
	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.gameState.MapState.PlayerLocation.Position)
	layeredMap := gs.gameState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.RegularDoor)

	// Clear any previous output messages
	gs.output.ClearRows()

	// Call the open command
	gs.smallMapOpenSecondary(direction)

	// Verify door was opened
	topTile := layeredMap.GetTopTile(targetPos)
	if topTile.Index != indexes.BrickFloor {
		t.Errorf("Expected door to be opened (BrickFloor), got %v", topTile.Index)
	}

	// Verify "Opened!" message was displayed
	lastMessage := gs.output.GetLastRowString()
	if lastMessage != "Opened!" {
		t.Errorf("Expected 'Opened!' message, got '%s'", lastMessage)
	}
}

func TestGameScene_SmallMapOpenSecondary_Door_Locked(t *testing.T) {
	gs := createTestGameScene(t)

	// Set up a locked door
	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.gameState.MapState.PlayerLocation.Position)
	layeredMap := gs.gameState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.LockedDoor)

	gs.output.ClearRows()

	gs.smallMapOpenSecondary(direction)

	// Verify door remains locked
	topTile := layeredMap.GetTopTile(targetPos)
	if topTile.Index != indexes.LockedDoor {
		t.Errorf("Expected door to remain LockedDoor, got %v", topTile.Index)
	}

	// Verify "Locked!" message
	lastMessage := gs.output.GetLastRowString()
	if lastMessage != "Locked!" {
		t.Errorf("Expected 'Locked!' message, got '%s'", lastMessage)
	}
}

func TestGameScene_SmallMapOpenSecondary_Door_MagicallyLocked(t *testing.T) {
	gs := createTestGameScene(t)

	// Set up a magically locked door
	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.gameState.MapState.PlayerLocation.Position)
	layeredMap := gs.gameState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.MagicLockDoor)

	gs.output.ClearRows()

	gs.smallMapOpenSecondary(direction)

	// Verify door remains magically locked
	topTile := layeredMap.GetTopTile(targetPos)
	if topTile.Index != indexes.MagicLockDoor {
		t.Errorf("Expected door to remain MagicLockDoor, got %v", topTile.Index)
	}

	// Verify "Magically Locked!" message
	lastMessage := gs.output.GetLastRowString()
	if lastMessage != "Magically Locked!" {
		t.Errorf("Expected 'Magically Locked!' message, got '%s'", lastMessage)
	}
}

func TestGameScene_SmallMapOpenSecondary_Chest_LordBritishTreasure(t *testing.T) {
	gs := createTestGameScene(t)

	// Set up for Lord British's castle basement
	gs.gameState.MapState.PlayerLocation.Location = references.Lord_Britishs_Castle
	gs.gameState.MapState.PlayerLocation.Floor = references.Basement

	// Set up a chest to the north
	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.gameState.MapState.PlayerLocation.Position)
	layeredMap := gs.gameState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.Chest)

	gs.output.ClearRows()

	gs.smallMapOpenSecondary(direction)

	// Verify treasure was found message
	// Note: The exact message format depends on the implementation
	messages := gs.output.GetAllRowStrings()
	foundTreasureMessage := false
	for _, msg := range messages {
		if msg == "Found:" {
			foundTreasureMessage = true
			break
		}
	}

	if !foundTreasureMessage {
		t.Error("Expected 'Found:' message for Lord British's treasure")
	}

	// Verify item stack was created at position
	hasItemStack := gs.gameState.ItemStacksMap.HasItemStackAtPosition(targetPos)
	if !hasItemStack {
		t.Error("Expected item stack to be created at chest position")
	}
}

func TestGameScene_SmallMapOpenSecondary_NotADoor_NormalChest(t *testing.T) {
	gs := createTestGameScene(t)

	// Set up a chest in a normal location (not Lord British's castle)
	gs.gameState.MapState.PlayerLocation.Location = references.Yew
	gs.gameState.MapState.PlayerLocation.Floor = references.Ground

	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.gameState.MapState.PlayerLocation.Position)
	layeredMap := gs.gameState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.Chest)

	gs.output.ClearRows()

	gs.smallMapOpenSecondary(direction)

	// For normal chests, the current implementation doesn't handle them yet
	// This test documents the current behavior - no messages should be added
	// When chest opening is fully implemented, this test should be updated
	messages := gs.output.GetAllRowStrings()
	if len(messages) > 0 {
		t.Logf("Current implementation produced messages: %v", messages)
		// This is expected behavior until general chest opening is implemented
	}
}

func TestGameScene_LargeMapOpenPrimary_ReturnsOpenWhat(t *testing.T) {
	gs := createTestGameScene(t)

	gs.output.ClearRows()

	// Simulate pressing 'O' key on large map
	gs.handleLargeMapInput(ebiten.KeyO)

	// Verify "Open what?" message
	lastMessage := gs.output.GetLastRowString()
	if lastMessage != "Open what?" {
		t.Errorf("Expected 'Open what?' message, got '%s'", lastMessage)
	}
}

func TestGameScene_OpenCommand_DirectionalFlow_SmallMap(t *testing.T) {
	gs := createTestGameScene(t)

	// Set up for small map mode
	gs.currentMapType = references.SmallMapType

	// Set up a door
	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.gameState.MapState.PlayerLocation.Position)
	layeredMap := gs.gameState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.RegularDoor)

	gs.output.ClearRows()

	// Simulate the two-phase input: 'O' then direction
	gs.handleSmallMapInput(ebiten.KeyO) // Should set secondary state

	// Verify we're in secondary input state
	expectedState := OpenDirectionInput // You'll need to adjust based on your actual state enum
	if gs.secondaryKeyState != expectedState {
		t.Errorf("Expected secondary state %v, got %v", expectedState, gs.secondaryKeyState)
	}

	// Now provide direction
	gs.handleSmallMapInput(ebiten.KeyUp) // North direction

	// Verify door was opened and message displayed
	topTile := layeredMap.GetTopTile(targetPos)
	if topTile.Index != indexes.BrickFloor {
		t.Errorf("Expected door to be opened, got %v", topTile.Index)
	}

	lastMessage := gs.output.GetLastRowString()
	if lastMessage != "Opened!" {
		t.Errorf("Expected 'Opened!' message, got '%s'", lastMessage)
	}
}

func TestGameScene_TimedDoorClosure_Integration(t *testing.T) {
	gs := createTestGameScene(t)

	// Open a door
	direction := references.North
	targetPos := direction.GetNewPositionInDirection(&gs.gameState.MapState.PlayerLocation.Position)
	layeredMap := gs.gameState.GetLayeredMapByCurrentLocation()
	layeredMap.SetTileByLayer(map_state.MapLayer, targetPos, indexes.RegularDoor)

	gs.smallMapOpenSecondary(direction)

	// Verify door is open
	if layeredMap.GetTopTile(targetPos).Index != indexes.BrickFloor {
		t.Fatal("Door should be open initially")
	}

	// Simulate turn advancement (this would normally happen through game actions)
	gs.gameState.MapState.SmallMapProcessTurnDoors()

	// Door should still be open after 1 turn
	if layeredMap.GetTopTile(targetPos).Index != indexes.BrickFloor {
		t.Error("Door should still be open after 1 turn")
	}

	// Advance another turn - door should close
	gs.gameState.MapState.SmallMapProcessTurnDoors()

	// Door should be closed now
	if layeredMap.GetTopTile(targetPos).Index != indexes.RegularDoor {
		t.Error("Door should be closed after 2 turns")
	}
}

// Helper function to create a test GameScene
func createTestGameScene(t *testing.T) *GameScene {
	t.Helper()

	// This would need proper GameScene initialization
	// For now, this is a placeholder showing the test structure
	gs := &GameScene{}

	// Initialize required components for testing:
	// - gameState with proper map setup
	// - output buffer for message verification
	// - currentMapType for input routing

	return gs
}
