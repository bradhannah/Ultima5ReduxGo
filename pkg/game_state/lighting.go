package game_state

import (
	"math"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type Lighting struct {
	turnsToExtinguishTorch int

	xTilesInMap, yTilesInMap int
	baselineFactor           float32
	baselineRadius           int
}

type DistanceMaskMap map[references.Position]int

func NewLighting(xTilesInMap, yTilesInMap int) Lighting {
	l := Lighting{}
	l.xTilesInMap = xTilesInMap
	l.yTilesInMap = yTilesInMap
	l.baselineFactor = 0.1
	l.baselineRadius = 1
	return l
}

func (l *Lighting) GetBaselineRadius() float32 {
	return l.baselineFactor
}

func (l *Lighting) HasTorchLit() bool {
	return l.turnsToExtinguishTorch > 0
}

func (l *Lighting) LightTorch() {
	l.turnsToExtinguishTorch = DefaultNumberOfTurnsUntilTorchExtinguishes
}

func (l *Lighting) AdvanceTurn() {
	l.turnsToExtinguishTorch = helpers.Max(l.turnsToExtinguishTorch-1, 0)
}

// ---------------------------------------------------------------------------
// BuildDistanceMap: world‑space, map‑based field‑of‑view
// ---------------------------------------------------------------------------

// BuildDistanceMap returns a map[MapPos]int where
//
//	value = 0 .. radius     → integer distance from the light centre
//	key   absent            → tile is completely dark
//
// You can map that distance to alpha any way you like, e.g.
//
//	alpha = 1 - float32(dist)/float32(radius)
func (l *Lighting) BuildDistanceMap(centre references.Position, factor float32) DistanceMaskMap {
	// 1.  Work out how far light reaches for this factor.
	radius := l.radiusForFactor(factor)
	radius2 := radius * radius // compare squared distance

	// 2.  Allocate roughly the right capacity so the map doesn't re‑hash.
	estCells := (radius*2 + 1) * (radius*2 + 1) // square that bounds the circle
	m := make(map[references.Position]int, estCells)

	// 3.  Walk the bounding square; keep only points inside the circle.
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy > radius2 { // outside circle → skip
				continue
			}

			worldX := centre.X + references.Coordinate(dx)
			worldY := centre.Y + references.Coordinate(dy)
			dist := int(math.Sqrt(float64(dx*dx + dy*dy))) // 0,1,2,..

			m[references.Position{X: worldX, Y: worldY}] = dist
		}
	}
	return m
}

// ---------------------------------------------------------------------------
// Helpers (unchanged logic, still receiver methods)
// ---------------------------------------------------------------------------

func (l *Lighting) radiusForFactor(factor float32) int {
	switch {
	case factor <= l.baselineFactor:
		return l.baselineRadius
	case factor >= 1:
		return l.maxRadius()
	default:
		norm := (factor - l.baselineFactor) / (1 - l.baselineFactor) // 0‥1
		return l.baselineRadius +
			int(math.Round(float64(norm)*float64(l.maxRadius()-l.baselineRadius)))
	}
}

func (l *Lighting) maxRadius() int {
	return int(math.Ceil(
		math.Hypot(float64(l.xTilesInMap-1)/2, float64(l.yTilesInMap-1)/2),
	))
}
