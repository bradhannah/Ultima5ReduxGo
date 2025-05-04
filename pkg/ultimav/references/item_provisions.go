package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

type Provision int

func (p Provision) ID() int {
	return int(p)
}

func (p Provision) Type() ItemType {
	return ItemTypeProvision
}

const NumberOfProvisions = 5

const NoProvision Provision = -1

const (
	Torches   Provision = 0
	Gem       Provision = 1
	Key       Provision = 2
	SkullKeys Provision = 3
	Food      Provision = 4
	Gold      Provision = 5 // not really a provision - but needed for descriptions
)

func (p Provision) GetSpriteIndex() indexes.SpriteIndex {
	switch p {
	case Food:
		return indexes.ItemFood
	case Gold:
		return indexes.ItemGold
	case Key:
		return indexes.ItemKey
	case Gem:
		return indexes.ItemGem
	case Torches:
		return indexes.ItemTorch
	}
	return indexes.Avatar_KeyIndex
}
