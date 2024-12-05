package references

import (
	_ "embed"
	"encoding/json"
	"log"
)

type ItemType string

const (
	ItemTypeReagent     ItemType = "Reagent"
	ItemTypeEquipment            = "Equipment"
	ItemTypeSpell                = "Spell"
	ItemTypeSpecialItem          = "Special"
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

	for _, equipment := range inventoryItems.inventoryItemsMap[ItemTypeEquipment] {
		inventoryItems.Equipment[Equipment(equipment.ItemIndex)] = equipment
	}

	for _, reagent := range inventoryItems.inventoryItemsMap[ItemTypeReagent] {
		inventoryItems.Reagent[Reagent(reagent.ItemIndex)] = reagent
	}

	for _, spell := range inventoryItems.inventoryItemsMap[ItemTypeSpell] {
		inventoryItems.Spell[Spell(spell.ItemIndex)] = spell
	}

	for _, scroll := range inventoryItems.inventoryItemsMap[ItemTypeScroll] {
		inventoryItems.Scroll[Scroll(scroll.ItemIndex)] = scroll
	}

	for _, special := range inventoryItems.inventoryItemsMap[ItemTypeSpecialItem] {
		inventoryItems.Special[SpecialItem(special.ItemIndex)] = special
	}

	for _, potion := range inventoryItems.inventoryItemsMap[ItemTypePotion] {
		inventoryItems.Potion[Potion(potion.ItemIndex)] = potion
	}

	for _, shard := range inventoryItems.inventoryItemsMap[ItemTypeShard] {
		inventoryItems.Shard[Shard(shard.ItemIndex)] = shard
	}

	for _, quest := range inventoryItems.inventoryItemsMap[ItemTypeQuestItem] {
		inventoryItems.QuestItem[QuestItem(quest.ItemIndex)] = quest
	}

	for _, moonstone := range inventoryItems.inventoryItemsMap[ItemTypeMoonstone] {
		inventoryItems.Moonstone[Moonstone(moonstone.ItemIndex)] = moonstone
	}

	for _, provision := range inventoryItems.inventoryItemsMap[ItemTypeProvision] {
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

func (i *InventoryItemReferences) GetEquipmentReference(equipment Equipment) InventoryItem {
	return InventoryItem{}
}
