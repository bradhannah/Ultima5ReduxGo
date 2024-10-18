package grammar

import "strings"

type TextCommand struct {
	Matches []Match
}

func NewTextCommand(matches []Match) *TextCommand {
	textCommand := TextCommand{}
	textCommand.Matches = matches

	return &textCommand
}

func (t *TextCommand) GetAutoComplete(command string) string {
	splitStr := strings.Split(command, " ")

	highIndex := len(splitStr) - 1
	bMatches, _ := t.Matches[highIndex].PartiallyMatches(splitStr[highIndex])
	if !bMatches {
		return command
	}

	return t.Matches[highIndex].GetSuffixHint(command)
}
