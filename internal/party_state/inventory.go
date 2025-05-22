package party_state

import (
	references2 "github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

const (
	MaxGold               = 9999
	MaxProvisionFood      = 9999
	MaxProvisionGems      = 99
	MaxProvisionTorches   = 99
	MaxProvisionKey       = 99
	MaxProvisionSkullKeys = 99
)

type ProvisionsQuantity struct {
	Food      uint16
	Gems      byte
	Torches   byte
	Keys      byte
	SkullKeys byte
}

type Inventory struct {
	Provisions ProvisionsQuantity
	Gold       uint16
}

func (i *Inventory) PutItemInInventory(item *references2.ItemAndQuantity) {
	if item.Item.Type() == references2.ItemTypeProvision {
		switch references2.Provision(item.Item.ID()) {
		case references2.Food:
			i.Provisions.Food = helpers.Min(i.Provisions.Food+uint16(item.Quantity), MaxProvisionFood)
		case references2.Key:
			i.Provisions.Keys = helpers.Min(i.Provisions.Keys+byte(item.Quantity), MaxProvisionKey)
		case references2.Gem:
			i.Provisions.Gems = helpers.Min(i.Provisions.Gems+byte(item.Quantity), MaxProvisionGems)
		case references2.Torches:
			i.Provisions.Torches = helpers.Min(i.Provisions.Torches+byte(item.Quantity), MaxProvisionTorches)
		case references2.SkullKeys:
			i.Provisions.SkullKeys = helpers.Min(i.Provisions.SkullKeys+byte(item.Quantity), MaxProvisionSkullKeys)
		case references2.Gold:
			i.Gold = helpers.Min(i.Gold+uint16(item.Quantity), MaxGold)
		default:
			panic("unhandled default case for PutItemInInventory")
		}
	}
}
