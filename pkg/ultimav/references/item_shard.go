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
	Hatred          = 1
	Cowardice       = 2
)
