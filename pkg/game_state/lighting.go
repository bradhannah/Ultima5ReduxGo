package game_state

import (
	"math"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const TorchTileDistance = 2

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
func (l *Lighting) BuildDistanceMap(
	centre references.Position,
	factor float32,
) DistanceMaskMap {

	// 1.  Convert factor → integer radius.
	radius := l.radiusForFactor(factor)

	//    Inclusive threshold: (r + 0.5)²   (makes the outline nicely round)
	rThresh := float64(radius) + 0.5
	r2 := rThresh * rThresh

	// 2.  Pre‑allocate close to the number of cells in the bounding square.
	est := (radius*2 + 1) * (radius*2 + 1)
	m := make(DistanceMaskMap, est)

	// 3.  Walk the bounding square; keep points inside the disc.
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {

			if float64(dx*dx+dy*dy) > r2 { // outside round disc → skip
				continue
			}

			world := references.Position{
				X: centre.X + references.Coordinate(dx),
				Y: centre.Y + references.Coordinate(dy),
			}

			dist := int(math.Sqrt(float64(dx*dx + dy*dy))) // 0‥radius
			m[world] = dist
		}
	}
	return m
}

func (l *Lighting) BuildLightSourceDistanceMap(lightSources LightSources) DistanceMaskMap {
	distanceMaskMap := make(DistanceMaskMap)

	// l.applyTorch(distanceMaskMap, references.Position{X: 80, Y: 103}, int(TorchTileDistance))
	for _, ls := range lightSources {
		l.applyTorch(distanceMaskMap, ls.Pos, ls.Tile.LightEmission)
	}

	//l.applyTorch(,
	return distanceMaskMap

}
func (l *Lighting) applyTorch(
	field map[references.Position]int,
	pos references.Position,
	radius int,
) {
	if radius <= 0 {
		return
	}

	// Pre‑compute the inclusive threshold: (r + 0.5)²
	rThresh := float64(radius) + 0.5
	r2 := rThresh * rThresh

	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			// Euclidean distance squared
			if float64(dx*dx+dy*dy) > r2 {
				continue // outside the nicely‑rounded disc
			}

			world := references.Position{
				X: pos.X + references.Coordinate(dx),
				Y: pos.Y + references.Coordinate(dy),
			}

			// Store integer Euclidean distance for your fade table
			dist := int(math.Sqrt(float64(dx*dx + dy*dy)))

			if cur, ok := field[world]; !ok || dist < cur {
				field[world] = dist // 0..radius
			}
		}
	}
}

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
