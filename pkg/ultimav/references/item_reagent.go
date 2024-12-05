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
	Ginseng              = 1
	Garlic               = 2
	SpiderSilk           = 3
	BloodMoss            = 4
	BlackPearl           = 5
	NightShade           = 6
	MandrakeRoot         = 7
)
