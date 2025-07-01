package party_state

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
)

type ProvisionsQuantity struct {
	Food      ItemQuantityLarge
	Gems      ItemQuantitySmall
	Torches   ItemQuantitySmall
	Keys      ItemQuantitySmall
	SkullKeys ItemQuantitySmall
}

type Inventory struct {
	Provisions ProvisionsQuantity
	Gold       ItemQuantityLarge

	Equipment    InventoryQuantities[references2.Equipment, *ItemQuantitySmall]
	Spells       InventoryQuantities[references2.Spell, *ItemQuantitySmall]
	Scrolls      InventoryQuantities[references2.Scroll, *ItemQuantityLarge]
	SpecialItems InventoryQuantities[references2.SpecialItem, *ItemQuantitySmall]
	QuestItems   InventoryQuantities[references2.QuestItem, *ItemQuantitySmall]
	Shards       InventoryQuantities[references2.Shard, *ItemQuantitySmall]
	Reagent      InventoryQuantities[references2.Reagent, *ItemQuantitySmall]
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.Equipment = NewInventoryQuantities[references2.Equipment, *ItemQuantitySmall]()
	inv.Spells = NewInventoryQuantities[references2.Spell, *ItemQuantitySmall]()
	inv.Scrolls = NewInventoryQuantities[references2.Scroll, *ItemQuantityLarge]()
	inv.SpecialItems = NewInventoryQuantities[references2.SpecialItem, *ItemQuantitySmall]()
	inv.QuestItems = NewInventoryQuantities[references2.QuestItem, *ItemQuantitySmall]()
	inv.Shards = NewInventoryQuantities[references2.Shard, *ItemQuantitySmall]()
	inv.Reagent = NewInventoryQuantities[references2.Reagent, *ItemQuantitySmall]()

	return inv
}

func (i *Inventory) PutItemInInventory(item *references2.ItemAndQuantity) {
	if item.Item.Type() == references2.ItemTypeProvision {
		switch references2.Provision(item.Item.ID()) {
		case references2.Food:
			i.Provisions.Food.IncrementBy(item.Quantity)
		case references2.Key:
			i.Provisions.Keys.IncrementBy(item.Quantity)
		case references2.Gem:
			i.Provisions.Gems.IncrementBy(item.Quantity)
		case references2.Torches:
			i.Provisions.Torches.IncrementBy(item.Quantity)
		case references2.SkullKeys:
			i.Provisions.SkullKeys.IncrementBy(item.Quantity)
		case references2.Gold:
			i.Gold.IncrementBy(item.Quantity)
		case references2.NoProvision:
			panic("unhandled default case with NoProvision")
		default:
			panic("unhandled default case for PutItemInInventory")
		}
	}
}
