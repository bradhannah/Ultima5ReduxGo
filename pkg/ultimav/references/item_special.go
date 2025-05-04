package references

type SpecialItem int

func (s SpecialItem) ID() int {
	return int(s)
}

func (s SpecialItem) Type() ItemType {
	return ItemTypeSpecialItem
}

const (
	Carpet      SpecialItem = 0
	Grapple     SpecialItem = 1
	Spyglass    SpecialItem = 2
	HMSCape     SpecialItem = 3
	PocketWatch SpecialItem = 4
	BlackBadge  SpecialItem = 5
	WoodenBox   SpecialItem = 6
	Sextant     SpecialItem = 7
)
