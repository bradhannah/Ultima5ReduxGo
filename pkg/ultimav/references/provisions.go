package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"

type Provision int

const NumberOfProvisions = 5

const NoProvision Provision = -1

const (
	Food = iota
	Gold
	Key
	Gem
	Torches
	SkullKeys
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
