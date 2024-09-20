package sprites

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image/png"
	"log"
	"time"
)

type SpriteAnimation struct {
	images                    []*ebiten.Image
	millisecondsBetweenChange int64
}

func NewSpriteAnimation(images []*ebiten.Image, millisecondsBetweenChange int64) *SpriteAnimation {
	return &SpriteAnimation{images: images, millisecondsBetweenChange: millisecondsBetweenChange}
}

func NewSpriteSlice(rawSprites [][]byte) []*ebiten.Image {
	images := make([]*ebiten.Image, len(rawSprites))
	for i, sprite := range rawSprites {
		img, err := png.Decode(bytes.NewReader(sprite))
		if err != nil {
			log.Fatal(err)
		}
		ebitenSprite := ebiten.NewImageFromImage(img)
		images[i] = ebitenSprite
	}
	return images
}

func (s *SpriteAnimation) GetCurrentSprite() *ebiten.Image {
	// get current time

	nFrame := int((time.Now().UnixMilli() / s.millisecondsBetweenChange) % int64(len(s.images)))
	return s.images[nFrame]
}
