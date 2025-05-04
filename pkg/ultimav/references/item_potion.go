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
	Yellow Potion = 1
	Red    Potion = 2
	Green  Potion = 3
	Orange Potion = 4
	Purple Potion = 5
	Black  Potion = 6
	White  Potion = 7
)
