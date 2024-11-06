package config

import "github.com/hajimehoshi/ebiten/v2"

type ScreenResolution struct {
	X int
	Y int
}

func GetWindowResolutionFromEbiten() ScreenResolution {
	width, height := ebiten.WindowSize()
	windowResolution := ScreenResolution{
		X: width,
		Y: height,
	}
	return windowResolution
}
