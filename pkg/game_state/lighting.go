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
	source references.Position,
	radius int,
	getTile references.GetTileFromPosition,
) {
	if radius <= 0 {
		return
	}

	roundR2 := math.Pow(float64(radius)+0.5, 2)

	br := l.gameDimensions.GetBottomRightWithoutOverflow()
	isWrapped := l.gameDimensions.IsWrappedMap()

	wrapPos := func(p references.Position) references.Position {
		if !isWrapped {
			return p
		}
		return *p.GetWrapped(br.X, br.Y)
	}
	inBounds := func(p references.Position) bool {
		return p.X >= 0 && p.Y >= 0 && p.X <= br.X && p.Y <= br.Y
	}

	// ---------- BFS without Chebyshev counter ----------
	q := []references.Position{source}
	seen := make(map[references.Position]struct{}, radius*radius)

	for len(q) > 0 {
		curPos := q[0]
		q = q[1:]

		if _, ok := seen[curPos]; ok {
			continue
		}
		seen[curPos] = struct{}{}

		if !isWrapped && !inBounds(curPos) {
			continue
		}

		dx := int(curPos.X - source.X)
		dy := int(curPos.Y - source.Y)
		if float64(dx*dx+dy*dy) > roundR2 { // outside round disc
			continue
		}

		euclid := int(math.Sqrt(float64(dx*dx + dy*dy)))
		if cur, ok := field[curPos]; !ok || euclid < cur {
			field[curPos] = euclid // light this tile (opaque or not)
		}

		w := wrapPos(curPos)
		tile := getTile(&w)

		// If tile is opaque, light adjacent walls but DO NOT enqueue further
		if tile.BlocksLight || tile.IsWindow {
			for ddy := -1; ddy <= 1; ddy++ {
				for ddx := -1; ddx <= 1; ddx++ {
					if ddx == 0 && ddy == 0 {
						continue
					}
					nb := references.Position{
						X: curPos.X + references.Coordinate(ddx),
						Y: curPos.Y + references.Coordinate(ddy),
					}
					nb = wrapPos(nb)
					if !isWrapped && !inBounds(nb) {
						continue
					}

					ndx := int(nb.X - source.X)
					ndy := int(nb.Y - source.Y)
					if float64(ndx*ndx+ndy*ndy) > roundR2 {
						continue
					}

					if getTile(&nb).IsWall() {
						distWall := int(math.Sqrt(float64(ndx*ndx + ndy*ndy)))
						if cur, ok := field[nb]; !ok || distWall < cur {
							field[nb] = distWall
						}
					}
				}
			}
			continue // stop flood at the opaque tile itself
		}

		// Transparent tile → enqueue 4‑way neighbours *if inside round disc*
		for _, dir := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nb := references.Position{
				X: curPos.X + references.Coordinate(dir[0]),
				Y: curPos.Y + references.Coordinate(dir[1]),
			}
			ndx := int(nb.X - source.X)
			ndy := int(nb.Y - source.Y)
			if float64(ndx*ndx+ndy*ndy) <= roundR2 {
				q = append(q, wrapPos(nb))
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
