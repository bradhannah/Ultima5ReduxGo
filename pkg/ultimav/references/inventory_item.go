package references

import (
	_ "embed"
	"encoding/json"
	"log"
)

type ItemType string

const (
	ItemTypeReagent     ItemType = "ItemTypeReagent"
	ItemTypeEquipment            = "Equipment"
	ItemTypeSpell                = "Spell"
	ItemTypeSpecialItem          = "SpecialItem"
	ItemTypeScroll               = "Scroll"
	ItemTypePotion               = "Potion"
	ItemTypeShard                = "Shard"
	ItemTypeQuestItem            = "QuestItem"
	ItemTypeMoonstone            = "Moonstone"
	ItemTypeProvision            = "Provision"
)

var (
	//go:embed data/InventoryDetails.json
	inventoryDetails []byte
)

type inventoryItemsMap map[ItemType][]InventoryItem

type InventoryItemReferences struct {
	inventoryItemsMap inventoryItemsMap

	Equipment map[Equipment]InventoryItem
}

func NewInventoryItemsReferences() *InventoryItemReferences {
	var inventoryItems InventoryItemReferences
	err := json.Unmarshal(inventoryDetails, &inventoryItems)
	if err != nil {
		log.Fatalf("error unmarshaling JSON: %v", err)
	}

	inventoryItems.Equipment = make(map[Equipment]InventoryItem)
	for _, equipment := range inventoryItems.inventoryItemsMap[ItemTypeEquipment] {
		inventoryItems.Equipment[Equipment(equipment.ItemIndex)] = equipment
	}

	return &inventoryItems
}

func (i *InventoryItemReferences) UnmarshalJSON(data []byte) error {
	var tempMap inventoryItemsMap
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	i.inventoryItemsMap = tempMap
	return nil
}

func (i *InventoryItemReferences) GetEquipmentReference(equipment Equipment) InventoryItem {
	return InventoryItem{}
}
