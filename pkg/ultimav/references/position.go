package references

import (
	"math"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

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

func (p *Position) GetFourDirectionsWrapped(maxX, maxY Coordinate) []Position {
	return []Position{
		{X: move(p.X, maxX, true, true), Y: p.Y},  // Left
		{X: move(p.X, maxX, false, true), Y: p.Y}, // Right
		{X: p.X, Y: move(p.Y, maxY, true, true)},  // Up
		{X: p.X, Y: move(p.Y, maxY, false, true)}, // Down
	}
}

func (p *Position) GetWrappedDistanceBetweenWrapped(p2 *Position, maxX, maxY Coordinate) float64 {
	dxRaw := float64(helpers.AbsInt(int(p.X - p2.X)))
	dyRaw := float64(helpers.AbsInt(int(p.Y - p2.Y)))
	dx := math.Min(dxRaw, float64(maxX)-dxRaw)
	dy := math.Min(dyRaw, float64(maxY)-dyRaw)
	return math.Sqrt(dx*dx + dy*dy)
}

func (p *Position) GetDistanceBetween(p2 *Position) float64 {
	dx := float64(p.X - p2.X)
	dy := float64(p.Y - p2.Y)
	preSqrt := dx*dx + dy*dy
	result := math.Sqrt(preSqrt)
	return result
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
func (p *Position) Equals(position *Position) bool {
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

func (p *Position) HeuristicTileDistance(b Position) int {
	return helpers.AbsInt(int(p.X-b.X)) + helpers.AbsInt(int(p.Y-b.Y))
}

func (p *Position) IsZeros() bool {
	return p.X == 0 && p.Y == 0
}

func (p *Position) IsNextTo(position Position) bool {
	return (p.X == position.X && helpers.AbsInt(int(p.Y-position.Y)) == 1) ||
		(p.Y == position.Y && helpers.AbsInt(int(p.X-position.X)) == 1)
}

// func (p *Position) GetSingleDirectionPositionCloserTo(position Position) *Position {
// 	if p.X < position.X {
// 		p.X++
// 	} else if p.X > position.X {
// 		p.X--
// 	} else if p.Y < position.Y {
// 		p.Y++
// 	} else if p.Y > position.Y {
// 		p.Y--
// 	}
// 	return p
// }
