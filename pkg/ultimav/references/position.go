package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"

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

func (p *Position) GetPositionToLeft() *Position {
	return &Position{p.X - 1, p.Y}
}

func (p *Position) GetPositionToRight() *Position {
	return &Position{p.X + 1, p.Y}
}

func (p *Position) GetPositionUp() *Position {
	return &Position{p.X, p.Y - 1}
}

func (p *Position) GetPositionDown() *Position {
	return &Position{p.X, p.Y + 1}
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

func (p *Position) GetWrapped(maxX Coordinate, maxY Coordinate) *Position {
	if p.X < 0 {
		p.X = p.X + maxX
	} else if p.X >= maxX {
		p.X = p.X % maxX
	}
	if p.Y < 0 {
		p.Y = p.Y + maxY
	} else if p.Y >= maxY {
		p.Y = p.Y % maxY
	}
	return p
}

func (p *Position) GetHash() int32 {
	const prime1, prime2 = 73856093, 19349663 // Use prime numbers for better distribution
	x := int32(p.X) * prime1
	y := int32(p.Y) * prime2
	return (x ^ y) & 0x7FFFFFFF // Ensures the result is within 32-bit signed int range
}

func (p *Position) IsWithinN(targetPosition *Position, n int) bool {
	return helpers.AbsInt(int(targetPosition.X-p.X)) <= n && helpers.AbsInt(int(targetPosition.Y-p.Y)) <= n
}

func (p *Position) Neighbors() []Position {
	return []Position{
		{X: p.X - 1, Y: p.Y}, // Left
		{X: p.X + 1, Y: p.Y}, // Right
		{X: p.X, Y: p.Y - 1}, // Up
		{X: p.X, Y: p.Y + 1}, // Down
	}
}
