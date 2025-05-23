package references

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

type ItemAndQuantity struct {
	Quantity int
	Item     Item
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
	items := make(ItemStacks, 0, total)

	const oneInXOddsOfGettingProvision = 3

	for i := 0; i < total; i++ {
		if helpers.OneInXOdds(oneInXOddsOfGettingProvision) {
			items = append(items, createRandomProvision())
		} else {
			item := ItemAndQuantity{}
			item.Item = getRandomNonSpecialEquipment()
			item.Quantity = 1
			// item.Item = NoProvision
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
	provision := Provision(helpers.RandomIntInRange(0, int(Torches)))
	item.Item = provision
	// item.Equipment = NoEquipment
	if provision == Gold {
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

func (i *ItemAndQuantity) GetAsProvision() Provision {
	return Provision(i.Item.ID())
}
