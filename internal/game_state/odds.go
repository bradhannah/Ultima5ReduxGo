package game_state

type TheOdds struct {
	oneInXLargeMapMonsterGeneration   int
	percentLikelyLargeMapMonsterMoves int
}

func NewDefaultTheOdds() TheOdds {
	theOdds := TheOdds{
		oneInXLargeMapMonsterGeneration:   32,
		percentLikelyLargeMapMonsterMoves: 75,
	}
	return theOdds
}

func (o *TheOdds) GetOneInXLargeMapMonsterGeneration() int {
	return o.oneInXLargeMapMonsterGeneration
}

func (o *TheOdds) SetGenerateLargeMapMonster(oneInX int) {
	o.oneInXLargeMapMonsterGeneration = oneInX
}

func (o *TheOdds) GetPercentLikeyLargeMapMonsterMoves() int {
	return o.percentLikelyLargeMapMonsterMoves
}

func (o *TheOdds) SetPercentLikeyLargeMapMonsterMoves(percent int) {
	o.percentLikelyLargeMapMonsterMoves = percent
}
