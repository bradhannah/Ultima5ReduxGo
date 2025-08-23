// Package main provides a utility to dump reference objects in JSON or YAML format.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// EnemyReferenceSafe is a dump-safe version of EnemyReference that avoids cyclic references
type EnemyReferenceSafe struct {
	KeyFrameTile         *references.Tile                `json:"key_frame_tile" yaml:"key_frame_tile"`
	Armour               int                             `json:"armour" yaml:"armour"`
	Damage               int                             `json:"damage" yaml:"damage"`
	Dexterity            int                             `json:"dexterity" yaml:"dexterity"`
	HitPoints            int                             `json:"hit_points" yaml:"hit_points"`
	Intelligence         int                             `json:"intelligence" yaml:"intelligence"`
	MaxPerMap            int                             `json:"max_per_map" yaml:"max_per_map"`
	Strength             int                             `json:"strength" yaml:"strength"`
	TreasureNumber       int                             `json:"treasure_number" yaml:"treasure_number"`
	EnemyAbilities       map[string]bool                 `json:"enemy_abilities" yaml:"enemy_abilities"`
	AdditionalEnemyFlags references.AdditionalEnemyFlags `json:"additional_enemy_flags" yaml:"additional_enemy_flags"`
	AttackRange          int                             `json:"attack_range" yaml:"attack_range"`
	// Friend omitted to avoid cyclic reference
}

// toFriendlyAbilityMap converts enemy abilities to a string-keyed map for JSON output
func toFriendlyAbilityMap(abilities map[references.EnemyAbility]bool) map[string]bool {
	friendly := make(map[string]bool)
	for k, v := range abilities {
		key := fmt.Sprintf("%d_%s", k, references.EnemyAbilityToString(k))
		friendly[key] = v
	}
	return friendly
}

// toSafeEnemyReferences converts EnemyReferences to a dump-safe format
func toSafeEnemyReferences(refs []references.EnemyReference) []EnemyReferenceSafe {
	safe := make([]EnemyReferenceSafe, len(refs))
	for i, e := range refs {
		safe[i] = EnemyReferenceSafe{
			KeyFrameTile:         e.KeyFrameTile,
			Armour:               e.Armour,
			Damage:               e.Damage,
			Dexterity:            e.Dexterity,
			HitPoints:            e.HitPoints,
			Intelligence:         e.Intelligence,
			MaxPerMap:            e.MaxPerMap,
			Strength:             e.Strength,
			TreasureNumber:       e.TreasureNumber,
			EnemyAbilities:       toFriendlyAbilityMap(e.EnemyAbilities),
			AdditionalEnemyFlags: e.AdditionalEnemyFlags,
			AttackRange:          e.AttackRange,
		}
	}
	return safe
}

// buildTalkDumpOutput converts raw talk scripts to a better structured format
func buildTalkDumpOutput(raw map[references.SmallMapMasterTypes][]*references.TalkScript) map[string][]map[string]interface{} {
	out := make(map[string][]map[string]interface{})
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
			entry["labels"] = ts.Labels
			npcDump = append(npcDump, entry)
		}
		out[fileName] = npcDump
	}
	return out
}

// dumpReferences handles the logic for dumping references to the specified output directory and format.
func dumpReferences(outputDir, outputFormat string) error {
	if outputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	if outputFormat != "json" && outputFormat != "yaml" {
		return fmt.Errorf("output format must be either 'json' or 'yaml'")
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputDir, 0750); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Log file creation
	logFile := fmt.Sprintf("%s/dump_log.txt", outputDir)
	logContent := fmt.Sprintf("Dump performed on: %s\nFormat: %s\n", time.Now().Format(time.RFC1123), outputFormat)
	if err := os.WriteFile(logFile, []byte(logContent), 0600); err != nil {
		return fmt.Errorf("failed to write log file: %w", err)
	}

	// Initialize GameReferences using the correct import path for NewUltimaVConfiguration
	gameConfig := config.NewUltimaVConfiguration()
	gameReferences, err := references.NewGameReferences(gameConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize game references: %w", err)
	}

	// Dump each field of GameReferences to its own file
	fields := map[string]interface{}{
		"LocationReferences":      gameReferences.LocationReferences,
		"DataOvl":                 gameReferences.DataOvl,
		"TileReferences":          gameReferences.TileReferences,
		"InventoryItemReferences": gameReferences.InventoryItemReferences,
		"LookReferences":          gameReferences.LookReferences.TileDescriptions, // Only dump TileDescriptions
		"NPCReferences":           gameReferences.NPCReferences,
		"DockReferences":          gameReferences.DockReferences,
		"EnemyReferences":         toSafeEnemyReferences(*gameReferences.EnemyReferences),
		"TalkReferences":          buildTalkDumpOutput(gameReferences.TalkReferences.GetTalkScripts()),
		// Note: Large map references excluded - too much data for reference purposes
	}

	for name, field := range fields {
		var fileExtension string
		if outputFormat == "yaml" {
			fileExtension = ".yaml"
		} else {
			fileExtension = ".json"
		}

		filePath := filepath.Join(outputDir, name+fileExtension)
		var output []byte

		if outputFormat == "yaml" {
			output, err = yaml.Marshal(field)
		} else {
			output, err = json.MarshalIndent(field, "", "  ")
		}

		if err != nil {
			return fmt.Errorf("failed to encode data for %s: %w", name, err)
		}
		if err := os.WriteFile(filePath, output, 0600); err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	return nil
}

func main() {
	var outputDir string
	var outputFormat string

	rootCmd := &cobra.Command{
		Use:   "ref_dump",
		Short: "Utility to dump reference objects in JSON or YAML format",
		RunE: func(_ *cobra.Command, _ []string) error {
			return dumpReferences(outputDir, outputFormat)
		},
	}

	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Directory to dump the files into (required)")
	rootCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", "Output format: 'json' or 'yaml'")
	if err := rootCmd.MarkFlagRequired("output"); err != nil {
		log.Printf("Error marking flag as required: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
