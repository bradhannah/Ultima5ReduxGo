package main

import (
	"fmt"
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
		grammar.StringMatch{
			Str:           "teleport",
			Description:   "Move to an X, Y coordinate on a given map",
			CaseSensitive: false,
		},
		grammar.IntMatch{IntMin: 0, IntMax: 255},
		grammar.IntMatch{IntMin: 0, IntMax: 255},
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
		grammar.StringMatch{
			Str:           "fy",
			Description:   "Go to the given floor",
			CaseSensitive: false,
		},
		grammar.IntMatch{IntMin: -1, IntMax: 5, Description: "Floor number"},
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
		grammar.StringMatch{
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
		grammar.StringMatch{
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
		grammar.StringMatch{
			Str:           "freemove",
			Description:   "Ignore boundaries when moving",
			CaseSensitive: false,
		}},
		func(s string, command *grammar.TextCommand) {
			d.gameScene.gameState.DebugOptions.FreeMove = !d.gameScene.gameState.DebugOptions.FreeMove
			d.dumpQuickState(fmt.Sprintf("FreeMove = %t", d.gameScene.gameState.DebugOptions.FreeMove))
		})
}
