package map_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type VisibilityCoords map[references.Coordinate]map[references.Coordinate]bool

func (v *VisibilityCoords) ResetVisibilityCoords(bDefault bool) {
	for row := range *v {
		for col := range (*v)[row] {
			(*v)[row][col] = bDefault
		}
	}
}

func (v *VisibilityCoords) SetVisibilityCoordsRectangle(topLeft *references.Position, bottomRight *references.Position, xMax references.Coordinate, yMax references.Coordinate, bWrap bool) {
	for x := topLeft.X; x < bottomRight.X; x++ {
		for y := topLeft.Y; y < bottomRight.Y; y++ {
			if bWrap {
				xy := (&references.Position{X: x, Y: y}).GetWrapped(xMax, yMax)
				(*v)[xy.X][xy.Y] = true
			} else {
				(*v)[x][y] = true
			}
		}
	}
}
