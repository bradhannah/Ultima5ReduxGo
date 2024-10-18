package grammar

import "strings"

type TextCommands []TextCommand

// OneOrMoreCommandsMatch
// Indicates if the string provide at least partially match the command
func (t *TextCommands) OneOrMoreCommandsMatch(command string) bool {
	return t.HowManyCommandsMatch(command) > 0
}

// HowManyCommandsMatch
// Returns the number of commands that currently match the text entered
func (t *TextCommands) HowManyCommandsMatch(command string) int {
	return len(*t.GetAllMatchingTextCommands(command))
}

func (t *TextCommands) GetAllMatchingTextCommands(command string) *TextCommands {
	splitCommand := strings.Split(command, " ")
	if splitCommand[len(splitCommand)-1] == "" {
		splitCommand = splitCommand[:len(splitCommand)-1]
	}

	allMatches := make(TextCommands, 0)

	nSplitCommands := len(splitCommand)
	for _, checkCommand := range *t {
		if nSplitCommands > len(checkCommand.Matches) {
			continue
		}
		for i := 0; i < nSplitCommands; i++ {
			bMatched, _ := checkCommand.Matches[i].PartiallyMatches(splitCommand[i])
			if bMatched && i == nSplitCommands-1 {
				allMatches = append(allMatches, checkCommand)
				break
			} else if !bMatched {
				break
			}
		}
	}

	return &allMatches
}
