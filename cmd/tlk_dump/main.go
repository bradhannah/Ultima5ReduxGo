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

func main() {
	format := flag.String("format", "json", "Output format: json or yaml")
	flag.Parse()

	// Load configuration
	cfg := config.NewUltimaVConfiguration()

	// Load DATA.OVL and all required game data
	dataOvl := references.NewDataOvl(cfg)

	// Load and parse all TLK files, expanding compressed words
	talkRefs := references.NewTalkReferences(cfg, dataOvl)

	// Dump all parsed TLK conversations to JSON
	outPath := "tlk_dump." + *format
	var outData []byte
	var err error

	if *format == "yaml" {
		outData, err = yaml.Marshal(talkRefs.GetTalkScripts())
	} else {
		outData, err = json.MarshalIndent(talkRefs.GetTalkScripts(), "", "  ")
	}

	if err != nil {
		fmt.Println("Error marshaling TLK data:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(outPath, outData, 0644); err != nil {
		fmt.Println("Error writing TLK dump:", err)
		os.Exit(1)
	}
	fmt.Printf("Dumped all parsed TLK conversations to: %s\n", outPath)
	fmt.Println("Instructions: Compare the output to https://wiki.ultimacodex.com/wiki/Ultima_V_transcript for validation.")
}
