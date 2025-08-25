package config

import (
	"github.com/bradhannah/Ultima5ReduxGo/internal/display"
	"github.com/hajimehoshi/ebiten/v2"
)

type ScreenResolution struct {
	X int
	Y int
}

// GetWindowResolutionFromEbiten returns current window resolution
// Deprecated: Use DisplayManager.GetCurrentSize() or GetResolution() instead
func GetWindowResolutionFromEbiten() ScreenResolution {
	if ebiten.IsFullscreen() {
		width, height := ebiten.Monitor().Size()
		return ScreenResolution{
			X: width,
			Y: height,
		}

	}
	width, height := ebiten.WindowSize()
	return ScreenResolution{
		X: width,
		Y: height,
	}
}

// SetWindowSize sets the window size using a global display manager
// This function is kept for backward compatibility
func SetWindowSize(screenResolution ScreenResolution) {
	// TODO: This should use a DisplayManager instance passed from config
	ebiten.SetWindowSize(screenResolution.X, screenResolution.Y)
	CenterWindow()
}

// CenterWindow centers the window on screen
// This function is kept for backward compatibility
func CenterWindow() {
	if ebiten.IsFullscreen() {
		return
	}

	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	screenResolution := GetWindowResolutionFromEbiten()

	ebiten.SetWindowPosition(monitorWidth/2-screenResolution.X/2, monitorHeight/2-screenResolution.Y/2)
}

// GetWindowResolutionFromDisplayManager gets resolution from a DisplayManager
func GetWindowResolutionFromDisplayManager(dm *display.DisplayManager) ScreenResolution {
	w, h := dm.GetCurrentSize()
	return ScreenResolution{X: w, Y: h}
}
