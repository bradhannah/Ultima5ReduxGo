package references

type Position struct {
	X, Y int16
}

func move(val, mapSize int16, decrease bool, bWrap bool) int16 {
	if decrease {
		if val == 0 {
			if bWrap {
				return mapSize - 1
			} else {
				return 0
			}
		}
		return val - 1
	}
	if val+1 >= mapSize {
		if bWrap {
			return 0
		} else {
			return val
		}
	}
	return val + 1
}

func (p *Position) GoLeft(bLarge bool) {
	mapSize := XSmallMapTiles
	if bLarge {
		mapSize = XLargeMapTiles
	}
	p.X = move(p.X, mapSize, true, bLarge)
}

func (p *Position) GoRight(bLarge bool) {
	mapSize := XSmallMapTiles
	if bLarge {
		mapSize = XLargeMapTiles
	}
	p.X = move(p.X, mapSize, false, bLarge)
}

func (p *Position) GoUp(bLarge bool) {
	mapSize := YSmallMapTiles
	if bLarge {
		mapSize = YLargeMapTiles
	}
	p.Y = move(p.Y, mapSize, true, bLarge)
}

func (p *Position) GoDown(bLarge bool) {
	mapSize := YSmallMapTiles
	if bLarge {
		mapSize = YLargeMapTiles
	}
	p.Y = move(p.Y, mapSize, false, bLarge)
}
func (p *Position) Equals(position Position) bool {
	return p.X == position.X && p.Y == position.Y
}
