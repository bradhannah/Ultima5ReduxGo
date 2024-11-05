package text

import "github.com/hajimehoshi/ebiten/v2"

const unscaledWindowWidth = 1920
const unscaledWindowHeight = 1080
const unscaledDefaultPointScale = 20

func GetScaleFactorByWindowSize() float64 {
	screenWidth, _ := ebiten.WindowSize()

	return float64(screenWidth) / float64(unscaledWindowWidth)
}

func GetScaledNumberToResolution(defaultFontPoint float64) float64 {
	return GetScaleFactorByWindowSize() * defaultFontPoint
}
