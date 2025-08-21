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

type DumpOutput map[string][]map[string]interface{}

func buildDumpOutput(raw map[references.SmallMapMasterTypes][]*references.TalkScript) DumpOutput {
	out := make(DumpOutput)
	for k, v := range raw {
		fileName := references.GetSmallMapTalkFile(k)
		var npcDump []map[string]interface{}
		for i, ts := range v {
			entry := make(map[string]interface{})
			entry["dialog_number"] = i
			var questionGroups []references.QuestionGroup
			if len(ts.Lines) > references.TalkScriptConstantsDescription {
				description := references.QuestionGroup{
					Options: []string{"description"},
					Script:  ts.Lines[references.TalkScriptConstantsDescription],
				}
				questionGroups = append(questionGroups, description)
			}
			if len(ts.Lines) > references.TalkScriptConstantsGreeting {
				greeting := references.QuestionGroup{
					Options: []string{"greeting"},
					Script:  ts.Lines[references.TalkScriptConstantsGreeting],
				}
				questionGroups = append(questionGroups, greeting)
			}
			questionGroups = append(questionGroups, ts.QuestionGroups...)
			entry["dialog"] = questionGroups
			npcDump = append(npcDump, entry)
		}
		out[fileName] = npcDump
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
