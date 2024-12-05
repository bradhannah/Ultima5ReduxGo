package references

import (
	_ "embed"
	"encoding/json"
	"log"
)

type ItemTypeStringIndex string

const (
	ItemTypeReagentStr     ItemTypeStringIndex = "Reagent"
	ItemTypeEquipmentStr                       = "Equipment"
	ItemTypeSpellStr                           = "Spell"
	ItemTypeSpecialItemStr                     = "Special"
	ItemTypeScrollStr                          = "Scroll"
	ItemTypePotionStr                          = "Potion"
	ItemTypeShardStr                           = "Shard"
	ItemTypeQuestItemStr                       = "QuestItem"
	ItemTypeMoonstoneStr                       = "Moonstone"
	ItemTypeProvisionStr                       = "Provision"
)

var (
	//go:embed data/InventoryDetails.json
	inventoryDetails []byte
)

type inventoryItemsMap map[ItemTypeStringIndex][]InventoryItem

type InventoryItemReferences struct {
	inventoryItemsMap inventoryItemsMap

	Equipment map[Equipment]InventoryItem
	Reagent   map[Reagent]InventoryItem
	Spell     map[Spell]InventoryItem
	Scroll    map[Scroll]InventoryItem
	Special   map[SpecialItem]InventoryItem
	Potion    map[Potion]InventoryItem
	Shard     map[Shard]InventoryItem
	QuestItem map[QuestItem]InventoryItem
	Moonstone map[Moonstone]InventoryItem
	Provision map[Provision]InventoryItem
}

func NewInventoryItemsReferences() *InventoryItemReferences {
	var inventoryItems InventoryItemReferences
	err := json.Unmarshal(inventoryDetails, &inventoryItems)
	if err != nil {
		log.Fatalf("error unmarshaling JSON: %v", err)
	}

	inventoryItems.Equipment = make(map[Equipment]InventoryItem)
	inventoryItems.Reagent = make(map[Reagent]InventoryItem)
	inventoryItems.Spell = make(map[Spell]InventoryItem)
	inventoryItems.Scroll = make(map[Scroll]InventoryItem)
	inventoryItems.Special = make(map[SpecialItem]InventoryItem)
	inventoryItems.Potion = make(map[Potion]InventoryItem)
	inventoryItems.Shard = make(map[Shard]InventoryItem)
	inventoryItems.QuestItem = make(map[QuestItem]InventoryItem)
	inventoryItems.Moonstone = make(map[Moonstone]InventoryItem)
	inventoryItems.Provision = make(map[Provision]InventoryItem)

	for _, equipment := range inventoryItems.inventoryItemsMap[ItemTypeEquipmentStr] {
		inventoryItems.Equipment[Equipment(equipment.ItemIndex)] = equipment
	}

	for _, reagent := range inventoryItems.inventoryItemsMap[ItemTypeReagentStr] {
		inventoryItems.Reagent[Reagent(reagent.ItemIndex)] = reagent
	}

	for _, spell := range inventoryItems.inventoryItemsMap[ItemTypeSpellStr] {
		inventoryItems.Spell[Spell(spell.ItemIndex)] = spell
	}

	for _, scroll := range inventoryItems.inventoryItemsMap[ItemTypeScrollStr] {
		inventoryItems.Scroll[Scroll(scroll.ItemIndex)] = scroll
	}

	for _, special := range inventoryItems.inventoryItemsMap[ItemTypeSpecialItemStr] {
		inventoryItems.Special[SpecialItem(special.ItemIndex)] = special
	}

	for _, potion := range inventoryItems.inventoryItemsMap[ItemTypePotionStr] {
		inventoryItems.Potion[Potion(potion.ItemIndex)] = potion
	}

	for _, shard := range inventoryItems.inventoryItemsMap[ItemTypeShardStr] {
		inventoryItems.Shard[Shard(shard.ItemIndex)] = shard
	}

	for _, quest := range inventoryItems.inventoryItemsMap[ItemTypeQuestItemStr] {
		inventoryItems.QuestItem[QuestItem(quest.ItemIndex)] = quest
	}

	for _, moonstone := range inventoryItems.inventoryItemsMap[ItemTypeMoonstoneStr] {
		inventoryItems.Moonstone[Moonstone(moonstone.ItemIndex)] = moonstone
	}

	for _, provision := range inventoryItems.inventoryItemsMap[ItemTypeProvisionStr] {
		inventoryItems.Provision[Provision(provision.ItemIndex)] = provision
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

// func (i *InventoryItemReferences) GetEquipmentReference() InventoryItem {
// 	return InventoryItem{}
// }

func (i *InventoryItemReferences) GetReferenceByItem(item Item) InventoryItem {
	switch item.Type() {
	case ItemTypeEquipment:
		return i.Equipment[Equipment(item.ID())]
	case ItemTypeReagent:
		return i.Reagent[Reagent(item.ID())]
	case ItemTypeSpell:
		return i.Spell[Spell(item.ID())]
	case ItemTypePotion:
		return i.Potion[Potion(item.ID())]
	case ItemTypeScroll:
		return i.Scroll[Scroll(item.ID())]
	case ItemTypeSpecialItem:
		return i.Special[SpecialItem(item.ID())]
	case ItemTypeMoonstone:
		return i.Moonstone[Moonstone(item.ID())]
	case ItemTypeProvision:
		return i.Provision[Provision(item.ID())]
	case ItemTypeQuestItem:
		return i.QuestItem[QuestItem(item.ID())]
	case ItemTypeShard:
		return i.Shard[Shard(item.ID())]
	}
	panic("Unexpected: item type invalid")
}
