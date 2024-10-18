package grammar

import "strings"

type TextCommand struct {
	Matches        []Match
	ExecuteCommand func(string, *[]Match)
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

func (t *TextCommand) GetTextCommandIfAtLeastPartialMatch(command string) *TextCommand {
	splitCommand := splitCommand(command)

	nSplitCommands := len(splitCommand)
	if nSplitCommands > len(t.Matches) {
		return nil
	}
	for i := 0; i < nSplitCommands; i++ {
		bMatched, _ := t.Matches[i].PartiallyMatches(splitCommand[i])
		if bMatched && i == nSplitCommands-1 {
			return t
		} else if !bMatched {
			return nil
		}
	}
	return nil
}

func (t *TextCommand) GetTextCommandIfPerfectMatch(command string) *TextCommand {
	splitCommand := splitCommand(command)

	if len(splitCommand) != len(t.Matches) {
		return nil
	}
	return t.GetTextCommandIfAtLeastPartialMatch(command)
}

func splitCommand(command string) []string {
	splitCommand := strings.Split(command, " ")
	if splitCommand[len(splitCommand)-1] == "" {
		splitCommand = splitCommand[:len(splitCommand)-1]
	}
	return splitCommand
}
