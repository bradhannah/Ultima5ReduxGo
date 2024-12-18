package references

type Moonstone int

func (m Moonstone) ID() int {
	return int(m)
}

func (m Moonstone) Type() ItemType {
	return ItemTypeMoonstone
}

const (
	NewMoon        Moonstone = 0
	CrescentWaxing           = 1
	FirstQuarter             = 2
	GibbousWaxing            = 3
	FullMoon                 = 4
	GibbousWaning            = 5
	LastQuarter              = 6
	CrescentWaning           = 7
)
