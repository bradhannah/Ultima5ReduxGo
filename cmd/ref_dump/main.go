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
	"gopkg.in/yaml.v2"
)

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

	// Dump references
	referencesData := map[string]interface{}{
		"locations": references.GetListOfAllLocationsWithDocksAsString(),
		// Add other references here as needed
	}

	var outputData []byte
	if outputFormat == "json" {
		outputData, err = json.MarshalIndent(referencesData, "", "  ")
	} else {
		outputData, err = yaml.Marshal(referencesData)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	outputFile := fmt.Sprintf("%s/references.%s", outputDir, outputFormat)
	if err := os.WriteFile(outputFile, outputData, 0600); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	// Dump each field of GameReferences to its own file
	fields := map[string]interface{}{
		"OverworldLargeMapReference":  gameReferences.OverworldLargeMapReference,
		"UnderworldLargeMapReference": gameReferences.UnderworldLargeMapReference,
		"LocationReferences":          gameReferences.LocationReferences,
		"DataOvl":                     gameReferences.DataOvl,
		"TileReferences":              gameReferences.TileReferences,
		"InventoryItemReferences":     gameReferences.InventoryItemReferences,
		"LookReferences":              gameReferences.LookReferences,
		"NPCReferences":               gameReferences.NPCReferences,
		"DockReferences":              *gameReferences.DockReferences,
		"EnemyReferences":             references.ToSafeEnemyReferences(*gameReferences.EnemyReferences),
		"TalkReferences":              gameReferences.TalkReferences,
	}

	// Debug: Print DockReferences before dumping
	fmt.Printf("DockReferences: %+v\n", gameReferences.DockReferences)
	// Debug: Print EnemyReferences after conversion
	fmt.Printf("EnemyReferences (safe): %+v\n", references.ToSafeEnemyReferences(*gameReferences.EnemyReferences))

	for name, field := range fields {
		filePath := filepath.Join(outputDir, name+".json")
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", filePath, err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(field); err != nil {
			return fmt.Errorf("failed to encode data for %s: %w", name, err)
		}
	}

	// Directly serialize safe enemy references for verification
	enemyReferencesSafe := references.ToSafeEnemyReferences(*gameReferences.EnemyReferences)
	enemyReferencesSafeFile := filepath.Join(outputDir, "EnemyReferencesSafe.json")
	enemyReferencesSafeF, err := os.Create(enemyReferencesSafeFile)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", enemyReferencesSafeFile, err)
	}
	encoder := json.NewEncoder(enemyReferencesSafeF)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(enemyReferencesSafe); err != nil {
		return fmt.Errorf("failed to encode data for EnemyReferencesSafe: %w", err)
	}
	enemyReferencesSafeF.Close()

	log.Printf("References dumped to %s", outputFile)
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
