package color

import "image/color"

//goland:noinspection GoUnusedGlobalVariable
var (
	Gray             = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xff}
	Black            = color.Black
	White            = color.White
	BlackSemi        = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x77}
	LighterBlackSemi = color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xAA}
	LighterBlueSemi  = color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0x77}
	Green            = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	Red              = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	Yellow           = color.RGBA{R: 0, G: 255, B: 255, A: 255}
	UltimaBlue       = color.RGBA{R: 22, G: 7, B: 178, A: 255}
)
