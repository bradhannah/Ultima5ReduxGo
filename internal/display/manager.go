package display

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// DisplayManager centralizes all screen/window management for the game.
// This is NOT an abstraction layer - it's a management layer that provides
// centralized control and change detection for display operations.
type DisplayManager struct {
	// Cached values to detect changes
	currentWidth  int
	currentHeight int
	isFullscreen  bool

	// Callbacks for resolution changes
	onResolutionChange []func(width, height int)
}

// NewDisplayManager creates a new display manager instance
func NewDisplayManager() *DisplayManager {
	w, h := ebiten.WindowSize()
	return &DisplayManager{
		currentWidth:  w,
		currentHeight: h,
		isFullscreen:  ebiten.IsFullscreen(),
	}
}

// Update checks for resolution changes and notifies listeners
// This should be called once per frame from the game loop
func (dm *DisplayManager) Update() {
	w, h := dm.GetCurrentSize()
	fs := ebiten.IsFullscreen()

	if w != dm.currentWidth || h != dm.currentHeight || fs != dm.isFullscreen {
		dm.currentWidth = w
		dm.currentHeight = h
		dm.isFullscreen = fs

		// Notify all registered listeners about the change
		for _, callback := range dm.onResolutionChange {
			callback(w, h)
		}
	}
}

// GetCurrentSize returns the current window or screen size
func (dm *DisplayManager) GetCurrentSize() (int, int) {
	if ebiten.IsFullscreen() {
		return ebiten.Monitor().Size()
	}
	return ebiten.WindowSize()
}

// SetWindowSize sets the window size (when not in fullscreen)
func (dm *DisplayManager) SetWindowSize(width, height int) {
	if !ebiten.IsFullscreen() {
		ebiten.SetWindowSize(width, height)
		dm.CenterWindow()
		dm.currentWidth = width
		dm.currentHeight = height
	}
}

// SetFullscreen toggles fullscreen mode
func (dm *DisplayManager) SetFullscreen(fullscreen bool) {
	ebiten.SetFullscreen(fullscreen)
	dm.isFullscreen = fullscreen
	// Resolution will be updated on next Update() call
}

// IsFullscreen returns the current fullscreen state
func (dm *DisplayManager) IsFullscreen() bool {
	return ebiten.IsFullscreen()
}

// CenterWindow centers the window on the monitor
func (dm *DisplayManager) CenterWindow() {
	if ebiten.IsFullscreen() {
		return
	}

	monitorW, monitorH := ebiten.Monitor().Size()
	windowW, windowH := ebiten.WindowSize()
	ebiten.SetWindowPosition(
		monitorW/2-windowW/2,
		monitorH/2-windowH/2,
	)
}

// RegisterResolutionChangeCallback adds a callback for resolution changes
func (dm *DisplayManager) RegisterResolutionChangeCallback(callback func(width, height int)) {
	dm.onResolutionChange = append(dm.onResolutionChange, callback)
}

// GetScaleFactor returns the scale factor based on a reference width
// This is useful for scaling UI elements based on screen size
func (dm *DisplayManager) GetScaleFactor(referenceWidth int) float64 {
	w, _ := dm.GetCurrentSize()
	return float64(w) / float64(referenceWidth)
}

// GetAspectRatio returns the current aspect ratio
func (dm *DisplayManager) GetAspectRatio() float64 {
	w, h := dm.GetCurrentSize()
	if h == 0 {
		return 1.0
	}
	return float64(w) / float64(h)
}
