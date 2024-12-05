package references

type Scroll int

func (s Scroll) ID() int {
	return int(s)
}

func (s Scroll) Type() ItemType {
	return ItemTypeScroll
}

const (
	ScrollVasLor     Scroll = 0
	ScrollRelHur            = 1
	ScrollInSanct           = 2
	ScrollInAn              = 3
	ScrollInQuasWis         = 4
	ScrollKalXenCorp        = 5
	ScrollInManiCorp        = 6
	ScrollAnTym             = 7
)
