package references

type Reagent int

func (r Reagent) ID() int {
	return int(r)
}

func (r Reagent) Type() ItemType {
	return ItemTypeReagent
}

const (
	SulfurAsh    Reagent = 0
	Ginseng      Reagent = 1
	Garlic       Reagent = 2
	SpiderSilk   Reagent = 3
	BloodMoss    Reagent = 4
	BlackPearl   Reagent = 5
	NightShade   Reagent = 6
	MandrakeRoot Reagent = 7
)
