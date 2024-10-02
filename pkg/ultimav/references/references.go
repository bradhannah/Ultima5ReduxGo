package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/config"

type GameReferences struct {
	OverworldLargeMapReference  *LargeMapReference
	UnderworldLargeMapReference *LargeMapReference

	SingleMapReferences *SingleMapReferences
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
	gameRefs.SingleMapReferences, err = NewSmallMapReferences(gameConfig)

	return gameRefs, nil
}
