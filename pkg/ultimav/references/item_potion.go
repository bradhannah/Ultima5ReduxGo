package references

type Potion int

func (p Potion) ID() int {
	return int(p)
}

func (p Potion) Type() ItemType {
	return ItemTypePotion
}

const (
	Blue   Potion = 0
	Yellow        = 1
	Red           = 2
	Green         = 3
	Orange        = 4
	Purple        = 5
	Black         = 6
	White         = 7
)
