package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"

type TheOdds struct {
	oneInXLargeMapMonsterGeneration int
}

func NewDefaultTheOdds() TheOdds {
	theOdds := TheOdds{
		oneInXLargeMapMonsterGeneration: 32,
	}
	return theOdds
}

func (o *TheOdds) GetLargeMapMonsterGeneration() int {
	return o.oneInXLargeMapMonsterGeneration
}

func (o *TheOdds) ShouldGenerateLargeMapMonster() bool {
	return helpers.OneInXOdds(o.oneInXLargeMapMonsterGeneration)
}

func (o *TheOdds) SetGenerateLargeMapMonster(oneInXLargeMapMonsterGeneration int) {
	o.oneInXLargeMapMonsterGeneration = oneInXLargeMapMonsterGeneration
}
