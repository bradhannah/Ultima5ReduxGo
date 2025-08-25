package display

// Resolution represents a screen resolution with width and height
type Resolution struct {
	Width  int
	Height int
}

// GetResolution returns the current resolution as a Resolution struct
func (dm *DisplayManager) GetResolution() Resolution {
	w, h := dm.GetCurrentSize()
	return Resolution{Width: w, Height: h}
}

// SetResolution sets the window to a specific resolution
func (dm *DisplayManager) SetResolution(res Resolution) {
	dm.SetWindowSize(res.Width, res.Height)
}
