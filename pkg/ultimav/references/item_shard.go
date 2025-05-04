package references

type Shard int

func (s Shard) ID() int {
	return int(s)
}

func (s Shard) Type() ItemType {
	return ItemTypeShard
}

const (
	Falsehood Shard = 0
	Hatred    Shard = 1
	Cowardice Shard = 2
)
