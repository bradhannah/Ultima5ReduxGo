package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/config"

type GameReferences struct {
	OverworldLargeMapReference  *LargeMapReference
	UnderworldLargeMapReference *LargeMapReference

	SingleMapReferences *SingleMapReferences
	DataOvl             *DataOvl
	TileReferences      *Tiles
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
	gameRefs.SingleMapReferences, err = NewSmallMapReferences(gameConfig, gameRefs.DataOvl)

	gameRefs.TileReferences = NewTileReferences()

	return gameRefs, nil
}
