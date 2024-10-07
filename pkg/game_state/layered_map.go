package game_state

import "github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"

func newLayeredMap(xMax int, yMax int, tileRefs *references.Tiles) *LayeredMap {
	const overflowTiles = 10
	layeredMap := LayeredMap{}
	layeredMap.tileRefs = tileRefs
	for i, _ := range layeredMap.Layers {
		layeredMap.Layers[i] = make(map[int]map[int]int)
		for j := -overflowTiles; j < yMax+10; j++ {
			layeredMap.Layers[i][j] = make(map[int]int)
		}
	}
	return &layeredMap
}

//func (l *LayeredMap) ClearTile(position *references.Position) {
//
//}

func (l *LayeredMap) GetTopTile(position *references.Position) *references.Tile {
	var tileValue int
	for i := totalLayers - 1; i >= 0; i-- {
		tileValue = l.Layers[i][int(position.X)][int(position.Y)]
		if tileValue <= 0 {
			continue
		}
		return l.tileRefs.GetTile(tileValue)
	}
	return nil
}
