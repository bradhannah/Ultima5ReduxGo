package text

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/display"
)

const (
	unscaledWindowWidth       = 1920
	unscaledWindowHeight      = 1080
	unscaledDefaultPointScale = 20
)

func GetScaleFactorByWindowSize(dm *display.DisplayManager) float64 {
	w, _ := dm.GetCurrentSize()
	return float64(w) / float64(unscaledWindowWidth)
}

func GetScaledNumberToResolution(dm *display.DisplayManager, defaultFontPoint float64) float64 {
	return GetScaleFactorByWindowSize(dm) * defaultFontPoint
}

// Deprecated: Legacy functions for backward compatibility with widgets
// These should be removed once widgets are updated to accept DisplayManager
func GetScaleFactorByWindowSizeLegacy() float64 {
	screenResolution := config.GetWindowResolutionFromEbiten()
	return float64(screenResolution.X) / float64(unscaledWindowWidth)
}

func GetScaledNumberToResolutionLegacy(defaultFontPoint float64) float64 {
	return GetScaleFactorByWindowSizeLegacy() * defaultFontPoint
}
