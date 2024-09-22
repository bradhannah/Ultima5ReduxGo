package sprites

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

// PercentXBasedPlacement
// Structure used for placing sprite on screen based on x and y positions based on percentage
// Note: there is no end percent Y as aspect ratio will always be maintained
type PercentXBasedPlacement struct {
	StartPercentX float64
	EndPercentX   float64
	StartPercentY float64
}

type PercentYBasedPlacement struct {
	StartPercentY float64
	EndPercentY   float64
}

// PercentBasedPlacement
// Structure used for placing sprite on screen based on x and y positions based on percentage
// Note: this will not maintain an aspect ratio and there is a bottom Y value
type PercentBasedPlacement struct {
	StartPercentX   float64
	EndXFontPercent float64
	StartPercentY   float64
	EndPercentY     float64
}

func GetRectangleFromPercents(placement PercentBasedPlacement) *image.Rectangle {
	screenWidth, screenHeight := ebiten.WindowSize()

	// get the x start and end values based on the percent
	var xLeft = float64(screenWidth) * placement.StartPercentX
	var xRight = float64(screenWidth) * placement.EndXFontPercent
	var yTop = float64(screenHeight) * placement.StartPercentY
	var yBottom = float64(screenHeight) * placement.EndPercentY

	return &image.Rectangle{
		Min: image.Point{X: int(xLeft), Y: int(yTop)},
		Max: image.Point{X: int(xRight), Y: int(yBottom)},
	}

}

func GetXSpriteWithPercents(rect image.Rectangle, placement PercentXBasedPlacement) *ebiten.DrawImageOptions {
	screenWidth, screenHeight := ebiten.WindowSize()

	op := &ebiten.DrawImageOptions{}

	// get the x start and end values based on the percent
	var xLeft = float64(screenWidth) * placement.StartPercentX
	var xRight = float64(screenWidth) * placement.EndPercentX
	var yTop = float64(screenHeight) * placement.StartPercentY

	targetWidth := xRight - xLeft
	originalImgWidth := rect.Bounds().Dx()

	// we scale to X value, and match the ratio on Y so the image doesn't stretch and history
	scaleX := targetWidth / float64(originalImgWidth)
	scaleY := scaleX
	op.GeoM.Scale(scaleX, scaleY)

	scaledWidth := scaleX * float64(originalImgWidth)

	topLeftX := (float64(screenWidth) - scaledWidth) / 2
	op.GeoM.Translate(topLeftX, yTop)

	return op
}

// GetYSpriteWithPercents
// Scales to the preferred Y % positions. It will center the X coordinate to the screen.
func GetYSpriteWithPercents(rect image.Rectangle, placement PercentYBasedPlacement) (*ebiten.DrawImageOptions, image.Point) {
	screenWidth, screenHeight := ebiten.WindowSize()

	op := &ebiten.DrawImageOptions{}

	// get the x start and end values based on the percent
	var yTop = float64(screenHeight) * placement.StartPercentY
	var yBottom = float64(screenHeight) * placement.EndPercentY

	targetHeight := yBottom - yTop
	originalImgHeight := rect.Bounds().Dy()

	// we scale to X value, and match the ratio on Y so the image doesn't stretch and history
	scaleY := targetHeight / float64(originalImgHeight)
	scaleX := scaleY
	op.GeoM.Scale(scaleX, scaleY)

	topLeftX := float64(screenWidth) / 2

	op.GeoM.Translate(topLeftX, yTop)

	return op, image.Point{
		X: int(topLeftX),
		Y: int(yTop),
	}

}
