package references

type Moonstone int

func (m Moonstone) ID() int {
	return int(m)
}

func (m Moonstone) Type() ItemType {
	return ItemTypeMoonstone
}

//goland:noinspection GoUnusedConst
const (
	NewMoon        Moonstone = 0
	CrescentWaxing Moonstone = 1
	FirstQuarter   Moonstone = 2
	GibbousWaxing  Moonstone = 3
	FullMoon       Moonstone = 4
	GibbousWaning  Moonstone = 5
	LastQuarter    Moonstone = 6
	CrescentWaning Moonstone = 7
)
