package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"

type ItemAndQuantity struct {
	Quantity  int
	Equipment Equipment
	Provision Provision
}

type Chest struct {
	Items []ItemAndQuantity
}

type ChestType int

const (
	LordBritishTreasure ChestType = iota
)

func getLordBritishItems(total int) Chest {
	var items []ItemAndQuantity

	const oneInXOddsOfGettingProvision = 3

	for i := 0; i < total; i++ {

		if helpers.OneInXOdds(oneInXOddsOfGettingProvision) {
			items = append(items, createRandomProvision())
		} else {
			item := ItemAndQuantity{}
			item.Equipment = getRandomNonSpecialEquipment()
			item.Quantity = 1
			item.Provision = NoProvision
		}
	}
	return Chest{
		Items: items,
	}
}

func getRandomNonSpecialEquipment() Equipment {
	const totalEquipmentIncludingSpecial = UusSanct + 1

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

func CreateNewChest(chestType ChestType) Chest {
	switch chestType {
	case LordBritishTreasure:
		const minTreasures, maxTreasures = 5, 20 // 15
		return getLordBritishItems(helpers.RandomIntInRange(minTreasures, maxTreasures))
	}
	return Chest{}
}
