package references

type TheOdds struct {
	oneInXMonsterGeneration int
}

func NewDefaultTheOdds() TheOdds {
	theOdds := TheOdds{
		oneInXMonsterGeneration: 32,
	}
	return theOdds
}

func (o *TheOdds) GetOneInXMonsterGeneration() int {
	return o.oneInXMonsterGeneration
}

func (o *TheOdds) SetMonsterGeneration(oneInX int) {
	o.oneInXMonsterGeneration = oneInX
}
