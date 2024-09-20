package sprites

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image/png"
	"log"
)

//go:embed assets/ultima.16.png
var ultima16Logo []byte // The image data will be stored in this byte slice

type IntroSprites struct {
	Ultima16Logo *ebiten.Image
}

func NewIntroSprites() *IntroSprites {
	introSprites := &IntroSprites{}
	img, err := png.Decode(bytes.NewReader(ultima16Logo))
	if err != nil {
		log.Fatal(err)
	}

	introSprites.Ultima16Logo = ebiten.NewImageFromImage(img)

	return introSprites
}
