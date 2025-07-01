package party_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type InventoryItemType interface {
	references.Equipment | references.Spell | references.Scroll | references.SpecialItem |
		references.QuestItem | references.Shard | references.Potion | references.Reagent
}

type InventoryQuantities[TK InventoryItemType, TV ItemQuantity] struct {
	quantities map[TK]*TV
}

func NewInventoryQuantities[TK InventoryItemType, TV ItemQuantity]() InventoryQuantities[TK, TV] {
	ei := InventoryQuantities[TK, TV]{}
	ei.quantities = make(map[TK]*TV)
	return ei
}

func (iq *InventoryQuantities[TK, TV]) GetQuantity(itemType TK) *TV {
	itemQuantity, ok := iq.quantities[itemType]
	if !ok {
		newItemQuantity := new(TV)
		iq.quantities[itemType] = newItemQuantity

		return newItemQuantity
	}

	return itemQuantity
}
