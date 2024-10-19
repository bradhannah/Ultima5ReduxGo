package sprites

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed assets/ultima.16.png
	ultima16Logo []byte // The image data will be stored in this byte slice
	//go:embed assets/ultima_fire.16.bmp
	flameLogo1 []byte
	//go:embed assets/ultima_fire2.16.bmp
	flameLogo2 []byte
	//go:embed assets/ultima_fire3.16.bmp
	flameLogo3 []byte
	//go:embed assets/ultima_fire4.16.bmp
	flameLogo4 []byte
	//go:embed assets/ReduxLogo.png
	reduxLogo []byte
)

type IntroSprites struct {
	Ultima16Logo   *ebiten.Image
	ReduxLogo      *ebiten.Image
	FlameAnimation *ImageAnimation
}

func NewIntroSprites() *IntroSprites {
	introSprites := &IntroSprites{}

	// static sprites
	introSprites.Ultima16Logo = NewPngSprite(ultima16Logo)
	introSprites.ReduxLogo = NewPngSprite(reduxLogo)

	// flame
	flameImages := NewSpriteSlice([][]byte{flameLogo1, flameLogo2, flameLogo3, flameLogo4})
	introSprites.FlameAnimation = NewImageAnimation(flameImages, 100)

	return introSprites
}
