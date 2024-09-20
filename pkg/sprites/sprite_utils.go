package sprites

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func NewPngSprite(rawBytes []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(rawBytes))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
