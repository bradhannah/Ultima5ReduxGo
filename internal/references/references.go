package references

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
)

type GameReferences struct {
	OverworldLargeMapReference  *LargeMapReference `json:"overworld_large_map_reference" yaml:"overworld_large_map_reference"`
	UnderworldLargeMapReference *LargeMapReference `json:"underworld_large_map_reference" yaml:"underworld_large_map_reference"`

	LocationReferences      *LocationReferences      `json:"location_references" yaml:"location_references"`
	DataOvl                 *DataOvl                 `json:"data_ovl" yaml:"data_ovl"`
	TileReferences          *Tiles                   `json:"tile_references" yaml:"tile_references"`
	InventoryItemReferences *InventoryItemReferences `json:"inventory_item_references" yaml:"inventory_item_references"`
	LookReferences          *LookReferences          `json:"look_references" yaml:"look_references"`
	NPCReferences           *NPCReferences           `json:"npc_references" yaml:"npc_references"`
	DockReferences          *DockReferences          `json:"dock_references" yaml:"dock_references"`
	EnemyReferences         *EnemyReferences         `json:"enemy_references" yaml:"enemy_references"`
	TalkReferences          *TalkReferences          `json:"talk_references" yaml:"talk_references"`
}

func NewGameReferences(gameConfig *config.UltimaVConfiguration) (*GameReferences, error) {
	gameRefs := &GameReferences{} //nolint:exhaustruct

	var err error
	gameRefs.OverworldLargeMapReference, err = NewLargeMapReference(gameConfig, OVERWORLD)
	if err != nil {
		return nil, err
	}

	gameRefs.UnderworldLargeMapReference, err = NewLargeMapReference(gameConfig, UNDERWORLD)
	if err != nil {
		return nil, err
	}

	gameRefs.DataOvl = NewDataOvl(gameConfig)
	gameRefs.LocationReferences, err = NewSmallMapReferences(gameConfig, gameRefs.DataOvl)
	if err != nil {
		log.Fatalf("Error when loading locations %e", err)
	}

	gameRefs.TileReferences = NewTileReferences()
	gameRefs.InventoryItemReferences = NewInventoryItemsReferences()
	gameRefs.LookReferences = NewLookReferences(gameConfig)

	gameRefs.NPCReferences = NewNPCReferences(gameConfig)
	gameRefs.DockReferences = NewDocks(gameConfig)

	gameRefs.EnemyReferences = NewAllEnemyReferences(gameConfig, gameRefs.TileReferences)

	gameRefs.TalkReferences = NewTalkReferences(gameConfig, gameRefs.DataOvl)

	return gameRefs, nil
}
