package references

type ItemStacksMap map[Position]ItemStack

func (i *ItemStacksMap) HasStackAtPosition(pos *Position) bool {
	_, exists := (*i)[*pos]
	return exists
}

func (i *ItemStacksMap) PopStackAtPosition(pos *Position) *ItemAndQuantity {
	oof := (*i)[*pos]
	item := oof.popTopItem()
	if !oof.HasItems() {
		delete(*i, *pos)
	}

	return &item
}

func (i *ItemStacksMap) PeekStackAtPosition(pos *Position) *ItemAndQuantity {
	oof, _ := (*i)[*pos]
	return oof.PeekTopItem()
}
