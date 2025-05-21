package references

type Scroll int

func (s Scroll) ID() int {
	return int(s)
}

func (s Scroll) Type() ItemType {
	return ItemTypeScroll
}

//goland:noinspection GoUnusedConst
const (
	ScrollVasLor     Scroll = 0
	ScrollRelHur     Scroll = 1
	ScrollInSanct    Scroll = 2
	ScrollInAn       Scroll = 3
	ScrollInQuasWis  Scroll = 4
	ScrollKalXenCorp Scroll = 5
	ScrollInManiCorp Scroll = 6
	ScrollAnTym      Scroll = 7
)
