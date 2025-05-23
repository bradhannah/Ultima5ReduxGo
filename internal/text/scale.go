package text

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
)

const (
	unscaledWindowWidth       = 1920
	unscaledWindowHeight      = 1080
	unscaledDefaultPointScale = 20
)

func GetScaleFactorByWindowSize() float64 {
	screenResolution := config.GetWindowResolutionFromEbiten()

	return float64(screenResolution.X) / float64(unscaledWindowWidth)
}

func GetScaledNumberToResolution(defaultFontPoint float64) float64 {
	return GetScaleFactorByWindowSize() * defaultFontPoint
}
