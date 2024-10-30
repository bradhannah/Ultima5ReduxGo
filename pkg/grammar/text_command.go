package grammar

import (
	"strconv"
	"strings"
)

const textCommandSeparator = " "

type ExecuteFunc func(string, *TextCommand)

type TextCommand struct {
	Matches        Matches
	ExecuteCommand ExecuteFunc
}

func NewTextCommand(matches Matches, executeFunc ExecuteFunc) *TextCommand {
	textCommand := TextCommand{}
	textCommand.Matches = matches
	textCommand.ExecuteCommand = executeFunc

	return &textCommand
}

func (t *TextCommand) GetIndexAsInt(nIndex int, command string) int {
	splitStr := strings.Split(command, textCommandSeparator)

	n, err := strconv.Atoi(splitStr[nIndex])
	if err != nil {
		return 0
	}
	return n
}

func (t *TextCommand) GetIndexAsString(nIndex int, command string) string {
	splitStr := strings.Split(strings.TrimSpace(command), textCommandSeparator)
	return splitStr[nIndex]
}

func (t *TextCommand) GetAutoComplete(command string) []string {
	splitStr := strings.Split(command, " ")

	highIndex := len(splitStr) - 1
	matches, _ := t.Matches[highIndex].GetPartialMatches(splitStr[highIndex])
	if len(matches) == 0 {
		return []string{""}
	}
	if len(matches) == 1 {
		return []string{t.Matches[highIndex].GetSuffixHint(splitStr[highIndex])}
	}

	return matches
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
