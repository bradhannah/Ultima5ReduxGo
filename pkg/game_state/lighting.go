package game_state

import (
	"math"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

const TorchTileDistance = 3

type Lighting struct {
	turnsToExtinguishTorch int
	gameDimensions         GameDimensions
	baselineFactor         float32
	baselineRadius         int
}

type DistanceMaskMap map[references.Position]int

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

func (l *Lighting) BuildLightSourceDistanceMap(
	lightSources LightSources,
	visibleFlags VisibilityCoords,
	avatarIgnitedTorch bool,
	centerPos references.Position,
	getTileFromPosition references.GetTileFromPosition) DistanceMaskMap {
	distanceMaskMap := make(DistanceMaskMap)

	if avatarIgnitedTorch {
		l.applyLightSource(distanceMaskMap, centerPos, TorchTileDistance, getTileFromPosition)
	}

	for _, ls := range lightSources {
		if visibleFlags[ls.Pos.X][ls.Pos.Y] {
			l.applyLightSource(distanceMaskMap, ls.Pos, ls.Tile.LightEmission, getTileFromPosition)
		}
	}

	return distanceMaskMap

}
func (l *Lighting) applyLightSource(
	field map[references.Position]int,
	pos references.Position,
	radius int,
	_ references.GetTileFromPosition,
) {
	if radius <= 0 {
		return
	}

	// Pre‑compute the inclusive threshold: (r + 0.5)²
	rThresh := float64(radius) + 0.5
	r2 := rThresh * rThresh

	//topLeft := l.gameDimensions.GetTopLeftExtent()
	bottomRight := l.gameDimensions.GetBottomRightWithoutOverflow()
	bIsWrapped := l.gameDimensions.IsWrappedMap()

	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			// Euclidean distance squared
			if float64(dx*dx+dy*dy) > r2 {
				continue // outside the nicely‑rounded disc
			}

			toPos := references.Position{
				X: pos.X + references.Coordinate(dx),
				Y: pos.Y + references.Coordinate(dy),
			}

			// we don't light the outer rim (overflow) tiles
			if !bIsWrapped {
				if toPos.X < 0 || toPos.Y < 0 || toPos.X > bottomRight.X || toPos.Y > bottomRight.Y {
					continue
				}
			} else {
				// is wrapped
				toPos = *toPos.GetWrapped(bottomRight.X, bottomRight.Y)
			}

			// TODO: still need to solve for the lightsource bleeding through the wall when there is a window that is not connected
			// to the light source
			// idea: can I just check for a window? I almost need to use a flood fill that ignores windows...
			// in the end I don't know mind if it bleeds out if it is in the viscinity of the light - but shouldn't show through if I am
			// not in close to a window with my light source

			// Store integer Euclidean distance for your fade table
			dist := int(math.Sqrt(float64(dx*dx + dy*dy)))

			// disabled: it appears that raycasting is not used in original game to detmermine if a light is obstructed
			// but I have kept it here "just in case"
			// if !pos.IsLineOfSight(toPos, getTileFromPosition) {
			// 	//continue
			// }

			if cur, ok := field[toPos]; !ok || dist < cur {
				field[toPos] = dist // 0..radius
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
