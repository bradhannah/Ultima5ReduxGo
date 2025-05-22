package sprites

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewPngSprite(rawBytes []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
