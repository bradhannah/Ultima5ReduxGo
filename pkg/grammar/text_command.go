package grammar

import "strings"

type TextCommands []TextCommand

type TextCommand struct {
	Matches []Match
}

func test() {

}

func NewTextCommand(matches []Match) *TextCommand {
	textCommand := TextCommand{}
	textCommand.Matches = matches

	return &textCommand
}

// OneOrMoreCommandsMatch
// Indicates if the string provide at least partially match the command
func (t *TextCommands) OneOrMoreCommandsMatch(command string) bool {
	return t.HowManyCommandsMatch(command) > 0
	//splitCommand := strings.Split(command, " ")
	//if splitCommand[len(splitCommand)-1] == "" {
	//	splitCommand = splitCommand[:len(splitCommand)-1]
	//}
	//nSplitCommands := len(splitCommand)
	//for _, checkCommand := range *t {
	//	if nSplitCommands > len(checkCommand.Matches) {
	//		continue
	//	}
	//	for i := 0; i < nSplitCommands; i++ {
	//		bMatched, _ := checkCommand.Matches[i].PartiallyMatches(splitCommand[i])
	//		if bMatched && i == nSplitCommands-1 {
	//			return true
	//		} else if !bMatched {
	//			break
	//		}
	//	}
	//}
	//return false
}

func (t *TextCommands) HowManyCommandsMatch(command string) int {
	splitCommand := strings.Split(command, " ")
	if splitCommand[len(splitCommand)-1] == "" {
		splitCommand = splitCommand[:len(splitCommand)-1]
	}

	nTotalMatches := 0

	nSplitCommands := len(splitCommand)
	for _, checkCommand := range *t {
		if nSplitCommands > len(checkCommand.Matches) {
			continue
		}
		for i := 0; i < nSplitCommands; i++ {
			bMatched, _ := checkCommand.Matches[i].PartiallyMatches(splitCommand[i])
			if bMatched && i == nSplitCommands-1 {
				nTotalMatches++
				break
			} else if !bMatched {
				break
			}
		}
	}
	return nTotalMatches
}
