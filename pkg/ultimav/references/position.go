package references

type Coordinate int16

type Position struct {
	X, Y Coordinate
}

func move(val, mapSize Coordinate, decrease bool, bWrap bool) Coordinate {
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

func (p *Position) GetPositionToLeft() Position {
	return Position{p.X - 1, p.Y}
}

func (p *Position) GetPositionToRight() Position {
	return Position{p.X + 1, p.Y}
}

func (p *Position) GetPositionUp() Position {
	return Position{p.X, p.Y - 1}
}

func (p *Position) GetPositionDown() Position {
	return Position{p.X, p.Y + 1}
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
