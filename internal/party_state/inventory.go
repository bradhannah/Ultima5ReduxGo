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
}

//const (
//	MaxGold               = 9999
//	MaxProvisionFood      = 9999
//	MaxProvisionGems      = 99
//	MaxProvisionTorches   = 99
//	MaxProvisionKey       = 99
//	MaxProvisionSkullKeys = 99
//)

func (i *Inventory) PutItemInInventory(item *references2.ItemAndQuantity) {
	if item.Item.Type() == references2.ItemTypeProvision {
		switch references2.Provision(item.Item.ID()) {
		case references2.Food:
			i.Provisions.Food.IncrementBy(uint16(item.Quantity))
		case references2.Key:
			i.Provisions.Keys.IncrementBy(uint16(item.Quantity))
		case references2.Gem:
			i.Provisions.Gems.IncrementBy(uint16(item.Quantity))
		case references2.Torches:
			i.Provisions.Torches.IncrementBy(uint16(item.Quantity))
		case references2.SkullKeys:
			i.Provisions.SkullKeys.IncrementBy(uint16(item.Quantity))
		case references2.Gold:
			i.Gold.IncrementBy(uint16(item.Quantity))
			// = helpers.Min(i.Gold+uint16(item.Quantity), MaxGold)
		default:
			panic("unhandled default case for PutItemInInventory")
		}
	}
}
