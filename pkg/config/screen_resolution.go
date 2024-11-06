package config

import "github.com/hajimehoshi/ebiten/v2"

type ScreenResolution struct {
	X int
	Y int
}

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

func SetWindowSize(screenResolution ScreenResolution) {
	ebiten.SetWindowSize(screenResolution.X, screenResolution.Y)
	CenterWindow()
}

func CenterWindow() {
	if ebiten.IsFullscreen() {
		return
	}

	monitorWidth, monitorHeight := ebiten.Monitor().Size()
	screenResolution := GetWindowResolutionFromEbiten()

	ebiten.SetWindowPosition(monitorWidth/2-screenResolution.X/2, monitorHeight/2-screenResolution.Y/2)
}
