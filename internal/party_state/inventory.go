package party_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
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

	Equipment    InventoryQuantities[references.Equipment, *ItemQuantitySmall]
	Spells       InventoryQuantities[references.Spell, *ItemQuantitySmall]
	Scrolls      InventoryQuantities[references.Scroll, *ItemQuantityLarge]
	SpecialItems InventoryQuantities[references.SpecialItem, *ItemQuantitySmall]
	QuestItems   InventoryQuantities[references.QuestItem, *ItemQuantitySmall]
	Shards       InventoryQuantities[references.Shard, *ItemQuantitySmall]
	Reagent      InventoryQuantities[references.Reagent, *ItemQuantitySmall]
}

func NewInventory() *Inventory {
	inv := &Inventory{}
	inv.Equipment = NewInventoryQuantities[references.Equipment, *ItemQuantitySmall]()
	inv.Spells = NewInventoryQuantities[references.Spell, *ItemQuantitySmall]()
	inv.Scrolls = NewInventoryQuantities[references.Scroll, *ItemQuantityLarge]()
	inv.SpecialItems = NewInventoryQuantities[references.SpecialItem, *ItemQuantitySmall]()
	inv.QuestItems = NewInventoryQuantities[references.QuestItem, *ItemQuantitySmall]()
	inv.Shards = NewInventoryQuantities[references.Shard, *ItemQuantitySmall]()
	inv.Reagent = NewInventoryQuantities[references.Reagent, *ItemQuantitySmall]()

	return inv
}

func (i *Inventory) PutItemInInventory(item *references.ItemAndQuantity) {
	if item.Item.Type() == references.ItemTypeProvision {
		switch references.Provision(item.Item.ID()) {
		case references.Food:
			i.Provisions.Food.IncrementBy(item.Quantity)
		case references.Key:
			i.Provisions.Keys.IncrementBy(item.Quantity)
		case references.Gem:
			i.Provisions.Gems.IncrementBy(item.Quantity)
		case references.Torches:
			i.Provisions.Torches.IncrementBy(item.Quantity)
		case references.SkullKeys:
			i.Provisions.SkullKeys.IncrementBy(item.Quantity)
		case references.Gold:
			i.Gold.IncrementBy(item.Quantity)
		case references.NoProvision:
			panic("unhandled default case with NoProvision")
		default:
			panic("unhandled default case for PutItemInInventory")
		}
	}
}
