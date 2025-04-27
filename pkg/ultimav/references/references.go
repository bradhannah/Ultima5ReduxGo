package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/config"

type GameReferences struct {
	OverworldLargeMapReference  *LargeMapReference
	UnderworldLargeMapReference *LargeMapReference

	LocationReferences      *LocationReferences
	DataOvl                 *DataOvl
	TileReferences          *Tiles
	InventoryItemReferences *InventoryItemReferences
	LookReferences          *LookReferences
	NPCReferences           *NPCReferences
	DockReferences          *DockReferences
	EnemyReferences         *EnemyReferences
}

func NewGameReferences(gameConfig *config.UltimaVConfiguration) (*GameReferences, error) {
	gameRefs := &GameReferences{}

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

	gameRefs.TileReferences = NewTileReferences()
	gameRefs.InventoryItemReferences = NewInventoryItemsReferences()
	gameRefs.LookReferences = NewLookReferences(gameConfig)

	gameRefs.NPCReferences = NewNPCReferences(gameConfig)
	gameRefs.DockReferences = NewDocks(gameConfig)

	gameRefs.EnemyReferences = NewAllEnemyReferences(gameConfig, gameRefs.TileReferences)

	return gameRefs, nil
}
