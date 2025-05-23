package references

// Item is an item that can be carried by a character.
type Item interface {
	ID() int        // Unique ID for the item
	Type() ItemType // Type of the item (e.g., Spell, Scroll, etc.)
}

// ItemType is the type of an item. An Item is a carriable object.
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
