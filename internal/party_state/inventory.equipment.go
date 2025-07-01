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

func (iq *InventoryQuantities[TK, TV]) IncrementByOne(itemType TK) {
	(*iq.GetQuantity(itemType)).IncrementByOne()
}

func (iq *InventoryQuantities[TK, TV]) DecrementByOne(itemType TK) bool {
	return (*iq.GetQuantity(itemType)).DecrementByOne()
}

func (iq *InventoryQuantities[TK, TV]) IncrementBy(itemType TK, incBy uint16) {
	(*iq.GetQuantity(itemType)).IncrementBy(incBy)
}

func (iq *InventoryQuantities[TK, TV]) DecrementBy(itemType TK, decBy uint16) bool {
	return (*iq.GetQuantity(itemType)).DecrementBy(decBy)
}

func (iq *InventoryQuantities[TK, TV]) Set(itemType TK, quantity uint16) {
	(*iq.GetQuantity(itemType)).Set(quantity)
}

func (iq *InventoryQuantities[TK, TV]) Get(itemType TK) uint16 {
	return (*iq.GetQuantity(itemType)).Get()
}

func (iq *InventoryQuantities[TK, TV]) HasSome(itemType TK) bool {
	return (*iq.GetQuantity(itemType)).HasSome()
}

func (iq *InventoryQuantities[TK, TV]) SetHasOne(itemType TK) {
	(*iq.GetQuantity(itemType)).SetHasOne()
}
