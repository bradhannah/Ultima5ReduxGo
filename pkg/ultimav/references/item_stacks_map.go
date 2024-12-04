package references

type ItemStacksMap struct {
	itemStacks map[Position]*ItemStack
}

func NewItemStacksMap() *ItemStacksMap {
	ism := &ItemStacksMap{}
	ism.itemStacks = make(map[Position]*ItemStack)
	return ism
}

func (i *ItemStacksMap) Push(pos *Position, itemStack *ItemStack) {
	if i.HasItemStackAtPosition(pos) {
		panic("Unexpected: pushing items on top of existing items")
	}
	i.itemStacks[*pos] = itemStack
}

func (i *ItemStacksMap) HasItemStackAtPosition(pos *Position) bool {
	_, exists := i.itemStacks[*pos]
	return exists
}

func (i *ItemStacksMap) Pop(pos *Position) *ItemAndQuantity {
	// oof := i.itemStacks[*pos]
	item := i.itemStacks[*pos].popTopItem()

	if !i.itemStacks[*pos].HasItems() {
		delete(i.itemStacks, *pos)
	}

	return &item
}

func (i *ItemStacksMap) Peek(pos *Position) *ItemAndQuantity {
	oof, _ := i.itemStacks[*pos]
	return oof.PeekTopItem()
}
