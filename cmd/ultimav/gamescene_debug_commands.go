package main

import (
	"fmt"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/internal/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/internal/map_units"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
)

func (d *DebugConsole) createDebugFunctions(gameScene *GameScene) *grammar.TextCommands {
	textCommands := make(grammar.TextCommands, 0)

	// Add each command by calling helper functions
	textCommands = append(textCommands, *d.createTeleportCommand(gameScene))
	textCommands = append(textCommands, *d.createFloorYCommand())
	textCommands = append(textCommands, *d.createFloorUpCommand())
	textCommands = append(textCommands, *d.createFloorDownCommand())
	textCommands = append(textCommands, *d.createFreeMoveCommand())
	textCommands = append(textCommands, *d.createTimeSet())
	textCommands = append(textCommands, *d.createQuickTime())
	textCommands = append(textCommands, *d.createGoSmall())
	textCommands = append(textCommands, *d.createResolutionUp())
	textCommands = append(textCommands, *d.createResolutionDown())
	textCommands = append(textCommands, *d.createFullScreenToggle())
	textCommands = append(textCommands, *d.createBuyBoat())
	textCommands = append(textCommands, *d.createToggleMonsterGen())
	textCommands = append(textCommands, *d.createMonsterGenerationOdds())
	textCommands = append(textCommands, *d.createRemoveAllMonsters())
	return &textCommands
}

func (d *DebugConsole) createRemoveAllMonsters() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "mon-delall",
			Description:   "Delete all monsters",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			d.gameScene.gameState.CurrentNPCAIController.RemoveAllEnemies()
			d.dumpQuickState("Removed all enemies")
		})
}

func (d *DebugConsole) createToggleMonsterGen() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "mon-toggle-gen",
			Description:   "Toggle monster generation",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			d.gameScene.gameState.DebugOptions.MonsterGen = !d.gameScene.gameState.DebugOptions.MonsterGen
			d.dumpQuickState(fmt.Sprintf("MonsterGen = %t", d.gameScene.gameState.DebugOptions.MonsterGen))
		})
}

func (d *DebugConsole) createMonsterGenerationOdds() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "mon-change-odds",
			Description:   "Set new monster generation odds",
			CaseSensitive: false,
		},
		grammar.MatchInt{IntMin: 0, IntMax: 1000, Description: "One in X odds of monster generation"},
	},
		func(s string, command *grammar.TextCommand) {
			outputStr := d.TextInput.GetText()
			n := command.GetIndexAsInt(1, outputStr)
			d.gameScene.gameState.TheOdds.SetGenerateLargeMapMonster(n)

			d.dumpQuickState(fmt.Sprintf("One In X odds of Monster=%d", n))
		})
}

func (d *DebugConsole) dumpQuickState(prefix string) {
	d.Output.AddRowStr(fmt.Sprintf("> %s\n  X=%d,Y=%d,Floor=%d",
		prefix,
		d.gameScene.gameState.MapState.PlayerLocation.Position.X,
		d.gameScene.gameState.MapState.PlayerLocation.Position.Y,
		d.gameScene.gameState.MapState.PlayerLocation.Floor))
}

// Helper function for the teleport command
func (d *DebugConsole) createTeleportCommand(gameScene *GameScene) *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "teleport",
			Description:   "Move to an X, Y coordinate on a given map",
			CaseSensitive: false,
		},
		grammar.MatchInt{IntMin: 0, IntMax: 255},
		grammar.MatchInt{IntMin: 0, IntMax: 255},
	}, func(s string, command *grammar.TextCommand) {
		outputStr := d.TextInput.GetText()
		gameScene.debugMoveOnMap(references.Position{
			X: references.Coordinate(command.GetIndexAsInt(1, outputStr)),
			Y: references.Coordinate(command.GetIndexAsInt(2, outputStr)),
		})
		d.dumpQuickState(fmt.Sprintf("Teleported to X=%d,Y=%d",
			int16(command.GetIndexAsInt(1, outputStr)),
			int16(command.GetIndexAsInt(2, outputStr)),
		))
	})
}

// Helper function for the floor command
func (d *DebugConsole) createFloorYCommand() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "fy",
			Description:   "Go to the given floor",
			CaseSensitive: false,
		},
		grammar.MatchInt{IntMin: -1, IntMax: 5, Description: "Floor number"},
	},
		func(s string, command *grammar.TextCommand) {
			outputStr := d.TextInput.GetText()

			res := d.gameScene.debugFloorY(references.FloorNumber(command.GetIndexAsInt(1, outputStr)))
			d.dumpQuickState(fmt.Sprintf("FloorTeleport Status=%t", res))
		})
}

// Helper function for the floor up command
func (d *DebugConsole) createFloorUpCommand() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "fu",
			Description:   "Teleport a floor up if one exists",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			res := d.gameScene.debugFloorUp()
			d.dumpQuickState(fmt.Sprintf("FloorUp Status=%t", res))
		})
}

// Helper function for the floor down command
func (d *DebugConsole) createFloorDownCommand() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "fd",
			Description:   "Teleport a floor down if one exists",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			res := d.gameScene.debugFloorDown()
			d.dumpQuickState(fmt.Sprintf("FloorDown Status=%t", res))
		})
}

func (d *DebugConsole) createFreeMoveCommand() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "freemove",
			Description:   "Ignore boundaries when moving",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			d.gameScene.gameState.DebugOptions.FreeMove = !d.gameScene.gameState.DebugOptions.FreeMove
			d.dumpQuickState(fmt.Sprintf("FreeMove = %t", d.gameScene.gameState.DebugOptions.FreeMove))
		})
}

func (d *DebugConsole) createTimeSet() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "tsh",
			Description:   "Set time to hour",
			CaseSensitive: false,
		},
		grammar.MatchInt{IntMin: 0, IntMax: 23, Description: "Hour"},
	},
		func(s string, command *grammar.TextCommand) {
			outputStr := d.TextInput.GetText()

			timeHour := byte(command.GetIndexAsInt(1, outputStr))
			d.gameScene.gameState.DateTime.Hour = timeHour
			d.dumpQuickState(fmt.Sprintf("Time Hour set=%b", timeHour))
		})
}

func (d *DebugConsole) createQuickTime() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "qt",
			Description:   "Quickly change time",
			CaseSensitive: false,
		},
		grammar.MatchStringList{
			Strings:       []string{"morning", "evening", "midnight", "noon", "dusk"},
			Description:   "General time of day",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			outputStr := strings.ToLower(d.TextInput.GetText())
			timeOfDay := command.GetIndexAsString(1, outputStr)
			switch timeOfDay {
			case "morning":
				d.gameScene.gameState.DateTime.SetTimeOfDay(datetime.Morning)
			case "evening":
				d.gameScene.gameState.DateTime.SetTimeOfDay(datetime.Evening)
			case "midnight":
				d.gameScene.gameState.DateTime.SetTimeOfDay(datetime.Midnight)
			case "noon":
				d.gameScene.gameState.DateTime.SetTimeOfDay(datetime.Noon)
			case "dusk":
				d.gameScene.gameState.DateTime.SetTimeOfDay(datetime.Dusk)
			}
			d.dumpQuickState(fmt.Sprintf("thing: %s", outputStr))
		})
}

func (d *DebugConsole) createGoSmall() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "gos",
			Description:   "Go to small map",
			CaseSensitive: false,
		},
		grammar.MatchStringList{
			Strings:       references.GetListOfAllSmallMaps(),
			Description:   "Small map locations",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			outputStr := strings.ToLower(d.TextInput.GetText())
			locationStr := command.GetIndexAsString(1, outputStr)
			slr := d.gameScene.gameReferences.LocationReferences.GetSmallLocationReference(locationStr)

			d.dumpQuickState(slr.FriendlyLocationName)
			d.gameScene.gameState.EnterBuilding(slr)
		})
}

func (d *DebugConsole) createBuyBoat() *grammar.TextCommand {
	docks := references.GetListOfAllLocationsWithDocksAsString()
	docks = append(docks, "avatar")

	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "buyboat",
			Description:   "Gets a frigate or skiff and puts at dock",
			CaseSensitive: false,
		},
		grammar.MatchStringList{
			Strings:              []string{"frigate", "skiff"},
			Description:          "Type of boat",
			CaseSensitive:        false,
			SingleCharacterInput: false,
		},
		grammar.MatchStringList{
			Strings:              docks,
			Description:          "Small map locations with docks",
			CaseSensitive:        false,
			SingleCharacterInput: false,
		},
	},
		func(_ string, command *grammar.TextCommand) {
			var dockPos references.Position
			var dockFloor references.FloorNumber //nolint:wsl

			outputStr := strings.ToLower(d.TextInput.GetText())
			boatType := command.GetIndexAsString(1, outputStr)
			locationStr := command.GetIndexAsString(2, outputStr)
			slr := d.gameScene.gameReferences.LocationReferences.GetSmallLocationReference(locationStr)

			if slr == nil {
				d.dumpQuickState("avatar")
				dockPos = d.gameScene.gameState.MapState.PlayerLocation.Position
				dockFloor = d.gameScene.gameState.MapState.PlayerLocation.Floor
			} else {
				d.dumpQuickState(slr.FriendlyLocationName)
				dockPos = d.gameScene.gameReferences.DockReferences.GetDockPositionByString(locationStr)
				dockFloor = 0
			}

			var boat map_units.NPCFriendly
			if strings.ToLower(boatType) == "frigate" {
				boat = *map_units.NewNPCFriendlyVehicle(references.FrigateVehicle, *references.NewNPCReferenceForVehicle(
					references.FrigateVehicle,
					dockPos,
					dockFloor,
				))
			} else {
				boat = *map_units.NewNPCFriendlyVehicle(references.SkiffVehicle, *references.NewNPCReferenceForVehicle(
					references.SkiffVehicle,
					dockPos,
					dockFloor,
				))
			}

			boat.NPCReference.Schedule.OverrideAllPositions(byte(dockPos.X), byte(dockPos.Y))
			boat.GetVehicleDetails().SetSkiffQuantity(1)

			bAddedVehicle := d.gameScene.gameState.LargeMapNPCAIController[references.OVERWORLD].GetNpcs().AddVehicle(boat)
			d.gameScene.gameState.FinishTurn()

			if !bAddedVehicle {
				d.dumpQuickState("Unable to add vehicle.")
			}
		})
}

// Helper function for the flor up command
func (d *DebugConsole) createResolutionUp() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "ru",
			Description:   "Increase the resolution",
			CaseSensitive: false,
		},
	},
		func(_ string, _ *grammar.TextCommand) {
			d.gameScene.gameConfig.IncrementHigherResolution()
			res := d.gameScene.gameConfig.GetCurrentTrackedWindowResolution()
			d.dumpQuickState(fmt.Sprintf("New Resolution: %dx%d", res.X, res.Y))
		})
}

func (d *DebugConsole) createResolutionDown() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "rd",
			Description:   "Shrink the resolution",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			d.gameScene.gameConfig.DecrementLowerResolution()
			res := d.gameScene.gameConfig.GetCurrentTrackedWindowResolution()
			d.dumpQuickState(fmt.Sprintf("New Resolution: %dx%d", res.X, res.Y))
		})
}

func (d *DebugConsole) createFullScreenToggle() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "fs",
			Description:   "Fullscreen Toggle",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			d.gameScene.gameConfig.SetFullScreen(!d.gameScene.gameConfig.SavedConfigData.FullScreen)
		})
}
