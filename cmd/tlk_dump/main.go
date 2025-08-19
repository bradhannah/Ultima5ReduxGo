package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"gopkg.in/yaml.v3"
)

type DumpOutput map[references.SmallMapMasterTypes][][]references.QuestionGroup

func buildDumpOutput(raw map[references.SmallMapMasterTypes][]*references.TalkScript) DumpOutput {
	out := make(DumpOutput)
	for k, v := range raw {
		var npcDump [][]references.QuestionGroup
		for _, ts := range v {
			// Include description and greeting as the first entries
			if len(ts.Lines) > references.TalkScriptConstantsDescription {
				description := references.QuestionGroup{
					Options: []string{"description"},
					Script:  ts.Lines[references.TalkScriptConstantsDescription],
				}
				npcDump = append(npcDump, []references.QuestionGroup{description})
			}
			if len(ts.Lines) > references.TalkScriptConstantsGreeting {
				greeting := references.QuestionGroup{
					Options: []string{"greeting"},
					Script:  ts.Lines[references.TalkScriptConstantsGreeting],
				}
				npcDump = append(npcDump, []references.QuestionGroup{greeting})
			}
			// Add the rest of the question groups
			npcDump = append(npcDump, ts.QuestionGroups)
		}
		out[k] = npcDump
	}
	return out
}

func main() {
	format := flag.String("format", "json", "Output format: json or yaml")
	flag.Parse()

	cfg := config.NewUltimaVConfiguration()
	dataOvl := references.NewDataOvl(cfg)
	talkRefs := references.NewTalkReferences(cfg, dataOvl)

	outPath := "tlk_dump." + *format
	var outData []byte
	var err error

	grouped := buildDumpOutput(talkRefs.GetTalkScripts())

	if *format == "yaml" {
		outData, err = yaml.Marshal(grouped)
	} else {
		outData, err = json.MarshalIndent(grouped, "", "  ")
	}

	if err != nil {
		fmt.Println("Error marshaling TLK data:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, outData, 0644); err != nil {
		fmt.Println("Error writing TLK dump:", err)
		os.Exit(1)
	}
	fmt.Printf("Dumped all grouped TLK questions to: %s\n", outPath)
	fmt.Println("Instructions: Compare the output to https://wiki.ultimacodex.com/wiki/Ultima_V_transcript for validation.")
}
