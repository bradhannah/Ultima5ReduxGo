package map_state

import (
	"math"

	"github.com/bradhannah/Ultima5ReduxGo/internal/references"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

const TorchTileDistance = 3
const DefaultNumberOfTurnsUntilTorchExtinguishes = 100

type Lighting struct {
	turnsToExtinguishTorch int
	gameDimensions         GameDimensions
	baselineFactor         float32
	baselineRadius         int
}

type DistanceMap map[references.Position]int

func NewLighting(gameDimensions GameDimensions) Lighting {
	l := Lighting{}
	l.gameDimensions = gameDimensions
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

func (l *Lighting) BuildGameScreenDistanceMap(centrePos references.Position) DistanceMap {
	dmm := make(DistanceMap)
	l.applyLightSource(dmm, centrePos, 0, 1)
	return dmm
}

func (l *Lighting) BuildLightSourceDistanceMap(
	lightSources LightSources,
	visibleFlags VisibilityCoords,
	avatarIgnitedTorch bool,
	centerPos references.Position) DistanceMap {
	distanceMaskMap := make(DistanceMap)

	if avatarIgnitedTorch {
		l.applyLightSource(distanceMaskMap, centerPos, TorchTileDistance, 1)
	}

	for _, ls := range lightSources {
		if visibleFlags[ls.Pos.X][ls.Pos.Y] {
			l.applyLightSource(distanceMaskMap, ls.Pos, ls.Tile.LightEmission, 1)
		}
	}

	return distanceMaskMap

}

func (l *Lighting) applyLightSource(
	distanceMap DistanceMap,
	centrePos references.Position,
	nRadius int, // 0 means “use factor”
	fFactor float32, // ignored when radius > 0
) {
	// Decide radius
	if nRadius <= 0 {
		nRadius = l.radiusForFactor(fFactor)
	}
	if nRadius <= 0 {
		return // nothing to do
	}

	// (r + 0.5)² → keeps outline round
	maxDist2 := math.Pow(float64(nRadius)+0.5, 2)

	br := l.gameDimensions.GetBottomRightWithoutOverflow()
	isWrapped := l.gameDimensions.IsWrappedMap()

	for dy := -nRadius; dy <= nRadius; dy++ {
		yy := dy * dy
		for dx := -nRadius; dx <= nRadius; dx++ {
			dist2 := float64(dx*dx + yy)
			if dist2 > maxDist2 { // outside disc
				continue
			}

			p := references.Position{
				X: centrePos.X + references.Coordinate(dx),
				Y: centrePos.Y + references.Coordinate(dy),
			}
			if isWrapped {
				p = *p.GetWrapped(br.X, br.Y)
			}

			d := int(math.Sqrt(dist2)) // Euclidean distance 0‥radius
			if cur := distanceMap[p]; cur == 0 || d < cur {
				distanceMap[p] = d
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
	xTilesVisibleOnGameScreen, yTilesVisibleOnGameScreen := l.gameDimensions.GetTilesVisibleOnScreen()
	return int(math.Ceil(
		math.Hypot(float64(xTilesVisibleOnGameScreen)/2, float64(yTilesVisibleOnGameScreen)/2),
	))
}
