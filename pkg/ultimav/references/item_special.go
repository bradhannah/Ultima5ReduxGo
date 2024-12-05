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
	Grapple                 = 1
	Spyglass                = 2
	HMSCape                 = 3
	PocketWatch             = 4
	BlackBadge              = 5
	WoodenBox               = 6
	Sextant                 = 7
)
