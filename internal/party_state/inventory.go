package party_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
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

func (i *Inventory) PutItemInInventory(item *references.ItemAndQuantity) {
	if item.Item.Type() == references.ItemTypeProvision {
		switch references.Provision(item.Item.ID()) {
		case references.Food:
			i.Provisions.Food = helpers.Min(i.Provisions.Food+uint16(item.Quantity), MaxProvisionFood)
		case references.Key:
			i.Provisions.Keys = helpers.Min(i.Provisions.Keys+byte(item.Quantity), MaxProvisionKey)
		case references.Gem:
			i.Provisions.Gems = helpers.Min(i.Provisions.Gems+byte(item.Quantity), MaxProvisionGems)
		case references.Torches:
			i.Provisions.Torches = helpers.Min(i.Provisions.Torches+byte(item.Quantity), MaxProvisionTorches)
		case references.SkullKeys:
			i.Provisions.SkullKeys = helpers.Min(i.Provisions.SkullKeys+byte(item.Quantity), MaxProvisionSkullKeys)
		case references.Gold:
			i.Gold = helpers.Min(i.Gold+uint16(item.Quantity), MaxGold)
		default:
			panic("unhandled default case for PutItemInInventory")
		}
	}
}
