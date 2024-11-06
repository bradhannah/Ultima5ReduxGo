package text

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
)

const unscaledWindowWidth = 1920
const unscaledWindowHeight = 1080
const unscaledDefaultPointScale = 20

func GetScaleFactorByWindowSize() float64 {
	screenResolution := config.GetWindowResolutionFromEbiten()

	return float64(screenResolution.X) / float64(unscaledWindowWidth)
}

func GetScaledNumberToResolution(defaultFontPoint float64) float64 {
	return GetScaleFactorByWindowSize() * defaultFontPoint
}
