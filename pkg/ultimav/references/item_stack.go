package references

import (
	"fmt"
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

type ItemAndQuantity struct {
	Quantity  int
	Equipment Equipment
	Provision Provision
}

type ItemStacks []ItemAndQuantity

type ItemStack struct {
	Items ItemStacks
}

type ItemStackType int

const (
	LordBritishTreasure ItemStackType = iota
)

func getLordBritishItems(total int) ItemStack {
	var items ItemStacks = make(ItemStacks, 0, total)

	const oneInXOddsOfGettingProvision = 3

	for i := 0; i < total; i++ {

		if helpers.OneInXOdds(oneInXOddsOfGettingProvision) {
			items = append(items, createRandomProvision())
		} else {
			item := ItemAndQuantity{}
			item.Equipment = getRandomNonSpecialEquipment()
			item.Quantity = 1
			item.Provision = NoProvision
			items = append(items, item)
		}
	}
	return ItemStack{
		Items: items,
	}
}

func getRandomNonSpecialEquipment() Equipment {
	const totalEquipmentIncludingSpecial = Ankh + 1

	for {
		equipment := Equipment(helpers.RandomIntInRange(0, int(totalEquipmentIncludingSpecial)))
		if equipment == ChaosSword || equipment == GlassSword || equipment == JeweledSword || equipment == MysticSword || equipment == MysticArmour {
			continue
		}
		return equipment
	}
}

func createRandomProvision() ItemAndQuantity {
	item := ItemAndQuantity{}
	item.Provision = Provision(helpers.RandomIntInRange(0, int(Torches)))
	item.Equipment = NoEquipment
	if item.Provision == Gold {
		const minGold, maxGold = 10, 60
		item.Quantity = helpers.RandomIntInRange(minGold, maxGold)
	} else {
		item.Quantity = helpers.RandomIntInRange(1, 2)
	}
	return item
}

func CreateNewItemStack(itemStackType ItemStackType) ItemStack {
	switch itemStackType {
	case LordBritishTreasure:
		const minTreasures, maxTreasures = 5, 20 // 15
		return getLordBritishItems(helpers.RandomIntInRange(minTreasures, maxTreasures))
	}
	return ItemStack{}
}

func (i *ItemStack) HasItems() bool {
	return len(i.Items) > 0
}

func (i *ItemStack) popTopItem() ItemAndQuantity {
	if len(i.Items) == 0 {
		log.Fatal("Can't pop from empty stack")
	}

	item := i.Items[len(i.Items)-1]
	i.Items = i.Items[:len(i.Items)-1]
	return item
}

func (i *ItemStack) PeekTopItem() *ItemAndQuantity {
	if !i.HasItems() {
		log.Fatal("Can't peek from empty stack")
	}
	return &i.Items[len(i.Items)-1]
}

func (i *ItemAndQuantity) GetFriendlyActionGetMessage() string {
	if i.Equipment != NoEquipment {
		return fmt.Sprintf("%s!", i.Equipment)
	}
	return ""
}
