package references

type Item interface {
	ID() int        // Unique ID for the item
	Type() ItemType // Type of the item (e.g., Spell, Scroll, etc.)
}

type ItemType int

const (
	ItemTypeReagent ItemType = iota
	ItemTypeEquipment
	ItemTypeSpell
	ItemTypeSpecialItem
	ItemTypeScroll
	ItemTypePotion
	ItemTypeShard
	ItemTypeQuestItem
	ItemTypeMoonstone
	ItemTypeProvision
)
