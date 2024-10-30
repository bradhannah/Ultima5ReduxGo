package grammar

import "fmt"

type TextCommands []TextCommand

func (t *TextCommands) GetNumOfCommands() int {
	return len(*t)
}

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
	allMatches := make(TextCommands, 0)
	for _, checkCommand := range *t {
		m := checkCommand.GetTextCommandIfAtLeastPartialMatch(command)
		if m == nil {
			continue
		}
		allMatches = append(allMatches, *m)
	}

	return &allMatches
}

func (t *TextCommands) GetAllPerfectMatches(command string) *TextCommands {
	allMatches := make(TextCommands, 0)
	for _, checkCommand := range *t {
		m := checkCommand.GetTextCommandIfPerfectMatch(command)
		if m == nil {
			continue
		}
		allMatches = append(allMatches, *m)
	}

	return &allMatches
}

func (t *TextCommands) GetFriendlyListOfAllCommands() string {
	var allCommands string

	for i, textCommand := range *t {
		if i != 0 {
			allCommands += ", "
		}
		allCommands += fmt.Sprintf("%s (%s)", textCommand.Matches[0].GetString(), textCommand.Matches[0].GetDescription())
	}
	return allCommands
}
