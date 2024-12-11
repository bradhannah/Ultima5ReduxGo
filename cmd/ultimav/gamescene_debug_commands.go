package main

import (
	"fmt"
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/datetime"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/grammar"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
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

	return &textCommands
}

func (d *DebugConsole) dumpQuickState(prefix string) {
	d.Output.AddRowStr(fmt.Sprintf("> %s\n  X=%d,Y=%d,Floor=%d",
		prefix,
		d.gameScene.gameState.Position.X,
		d.gameScene.gameState.Position.Y,
		d.gameScene.gameState.Floor))
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
		gameScene.DebugMoveOnMap(references.Position{
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

			res := d.gameScene.DebugFloorY(references.FloorNumber(command.GetIndexAsInt(1, outputStr)))
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
		}},
		func(s string, command *grammar.TextCommand) {
			res := d.gameScene.DebugFloorUp()
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
		}},
		func(s string, command *grammar.TextCommand) {
			res := d.gameScene.DebugFloorDown()
			d.dumpQuickState(fmt.Sprintf("FloorDown Status=%t", res))
		})
}

func (d *DebugConsole) createFreeMoveCommand() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "freemove",
			Description:   "Ignore boundaries when moving",
			CaseSensitive: false,
		}},
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
			Strings:       []string{"morning", "evening", "midnight", "noon"},
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
			d.gameScene.gameState.EnterBuilding(
				slr,
				d.gameScene.gameReferences.TileReferences,
			)
		})
}

func (d *DebugConsole) createBuyBoat() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "buyboat",
			Description:   "Gets a frigate or skiff and puts at dock",
			CaseSensitive: false,
		},
		grammar.MatchStringList{
			Strings:       []string{"frigate", "skiff"},
			Description:   "Type of boat",
			CaseSensitive: false,
		},
		grammar.MatchStringList{
			Strings:       references.GetListOfAllLocationsWithDocksAsString(),
			Description:   "Small map locations with docks",
			CaseSensitive: false,
		},
	},
		func(s string, command *grammar.TextCommand) {
			outputStr := strings.ToLower(d.TextInput.GetText())
			boatType := command.GetIndexAsString(1, outputStr)
			locationStr := command.GetIndexAsString(2, outputStr)
			slr := d.gameScene.gameReferences.LocationReferences.GetSmallLocationReference(locationStr)

			d.dumpQuickState(slr.FriendlyLocationName)

			// put the boat
			_ = boatType

		})
}

// Helper function for the flor up command
func (d *DebugConsole) createResolutionUp() *grammar.TextCommand {
	return grammar.NewTextCommand([]grammar.Match{
		grammar.MatchString{
			Str:           "ru",
			Description:   "Increase the resolution",
			CaseSensitive: false,
		}},
		func(s string, command *grammar.TextCommand) {
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
		}},
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
		}},
		func(s string, command *grammar.TextCommand) {
			d.gameScene.gameConfig.SetFullScreen(!d.gameScene.gameConfig.SavedConfigData.FullScreen)
		})
}
