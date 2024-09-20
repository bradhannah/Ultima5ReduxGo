package sprites

import "github.com/hajimehoshi/ebiten/v2"

func GetCenteredXSprite(image *ebiten.Image, nPixelsPerSide int, y int) *ebiten.DrawImageOptions {
	screenWidth, _ := ebiten.WindowSize()

	op := &ebiten.DrawImageOptions{}

	var xDiff = nPixelsPerSide * 2
	targetWidth := float64(screenWidth - xDiff)
	imgWidth := image.Bounds().Dx()

	scaleX := targetWidth / float64(imgWidth)
	scaleY := scaleX
	op.GeoM.Scale(scaleX, scaleY)

	scaledWidth := scaleX * float64(imgWidth)

	topLeftX := (float64(screenWidth) - scaledWidth) / 2
	op.GeoM.Translate(topLeftX, float64(y))

	return op
}
