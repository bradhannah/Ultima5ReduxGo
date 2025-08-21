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

	// Dump each field of GameReferences to its own file
	fields := map[string]interface{}{
		"LocationReferences": gameReferences.LocationReferences,
		"DataOvl":            gameReferences.DataOvl,
		"TileReferences":     gameReferences.TileReferences,
		"LookReferences":     gameReferences.LookReferences.TileDescriptions, // Only dump TileDescriptions
		"NPCReferences":      gameReferences.NPCReferences,
		"DockReferences":     gameReferences.DockReferences,
		"EnemyReferences":    references.ToSafeEnemyReferences(*gameReferences.EnemyReferences),
		"TalkReferences":     gameReferences.TalkReferences,
	}

	for name, field := range fields {
		filePath := filepath.Join(outputDir, name+".json")
		var output []byte

		output, err = json.MarshalIndent(field, "", "  ")

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
