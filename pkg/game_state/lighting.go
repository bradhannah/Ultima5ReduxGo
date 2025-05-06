package game_state

import (
	"math"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
)

type Lighting struct {
	TurnsToExtinguishTorch int

	xTilesInMap, yTilesInMap int
	baselineFactor           float32
	baselineRadius           int
}

type DistanceMaskMap [][]int

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
	return l.TurnsToExtinguishTorch > 0
}

func (l *Lighting) LightTorch() {
	l.TurnsToExtinguishTorch = DefaultNumberOfTurnsUntilTorchExtinguishes
}

// BuildMask returns a y×x 2‑D slice telling you which tiles are lit
// (`true`) or hidden (`false`).
//
//	centre – light source in *screen* coordinates
//	factor – 0.10 → 1.00 from GetVisibilityFactorWithoutTorch().
func (l *Lighting) BuildDistanceMask(centre references.Position, factor float32) DistanceMaskMap {
	w, h := l.xTilesInMap, l.yTilesInMap

	// Allocate & preset to -1 (dark)
	mask := make(DistanceMaskMap, h)
	for y := range mask {
		row := make([]int, w)
		for x := range row {
			row[x] = -1
		}
		mask[y] = row
	}

	radius := l.radiusForFactor(factor)
	radius2 := radius * radius // compare squared distances

	for y := 0; y < h; y++ {
		dy := y - int(centre.Y)
		for x := 0; x < w; x++ {
			dx := x - int(centre.X)
			d2 := dx*dx + dy*dy
			if d2 <= radius2 { // inside lit circle
				mask[y][x] = int(math.Sqrt(float64(d2))) // 0,1,2…
			}
		}
	}
	return mask
}

// ─── Helpers (receiver methods) ────────────────────────────────────────────────
func (l *Lighting) radiusForFactor(factor float32) int {
	if factor <= l.baselineFactor {
		return l.baselineRadius
	}
	if factor >= 1 {
		return l.maxRadius()
	}
	norm := (factor - l.baselineFactor) / (1 - l.baselineFactor) // 0‥1
	return l.baselineRadius +
		int(math.Round(float64(norm)*float64(l.maxRadius()-l.baselineRadius)))
}

func (l *Lighting) maxRadius() int {
	return int(math.Ceil(
		math.Hypot(float64(l.xTilesInMap-1)/2, float64(l.yTilesInMap-1)/2),
	))
}

//func (d *DistanceMaskMap) ApplyMaskTo
