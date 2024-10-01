package ultimav

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

type Position struct {
	X, Y int16
}

func move(val, mapSize int16, decrease bool) int16 {
	if decrease {
		if val == 0 {
			return mapSize - 1
		}
		return val - 1
	}
	if val+1 >= mapSize {
		return 0
	}
	return val + 1
}

func (p *Position) GoLeft(bLarge bool) {
	mapSize := references.XSmallMapTiles
	if bLarge {
		mapSize = references.XLargeMapTiles
	}
	p.X = move(p.X, mapSize, true)
}

func (p *Position) GoRight(bLarge bool) {
	mapSize := references.XSmallMapTiles
	if bLarge {
		mapSize = references.XLargeMapTiles
	}
	p.X = move(p.X, mapSize, false)
}

func (p *Position) GoUp(bLarge bool) {
	mapSize := references.YSmallMapTiles
	if bLarge {
		mapSize = references.YLargeMapTiles
	}
	p.Y = move(p.Y, mapSize, true)
}

func (p *Position) GoDown(bLarge bool) {
	mapSize := references.YSmallMapTiles
	if bLarge {
		mapSize = references.YLargeMapTiles
	}
	p.Y = move(p.Y, mapSize, false)
}
