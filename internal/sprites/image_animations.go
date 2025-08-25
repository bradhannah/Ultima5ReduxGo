package sprites

import (
	"bytes"
	"image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type ImageAnimation struct {
	images                    []*ebiten.Image
	millisecondsBetweenChange int64
}

func NewImageAnimation(images []*ebiten.Image, millisecondsBetweenChange int64) *ImageAnimation {
	return &ImageAnimation{images: images, millisecondsBetweenChange: millisecondsBetweenChange}
}

func NewSpriteSlice(rawSprites [][]byte) []*ebiten.Image {
	images := make([]*ebiten.Image, len(rawSprites))
	for i, sprite := range rawSprites {
		img, err := png.Decode(bytes.NewReader(sprite))
		if err != nil {
			log.Fatal(err) // Embedded animation data should always load successfully
		}
		ebitenSprite := ebiten.NewImageFromImage(img)
		images[i] = ebitenSprite
	}
	return images
}

func (s *ImageAnimation) GetCurrentImage() *ebiten.Image {
	// get current time - NOTE: This is non-deterministic, prefer GetCurrentImageAtTime() for core logic

	nFrame := int((time.Now().UnixMilli() / s.millisecondsBetweenChange) % int64(len(s.images)))
	return s.images[nFrame]
}

// GetCurrentImageAtTime returns the animation frame for a specific time (deterministic)
// Use this variant for core game logic that needs to be reproducible
func (s *ImageAnimation) GetCurrentImageAtTime(elapsedMs int64) *ebiten.Image {
	nFrame := int((elapsedMs / s.millisecondsBetweenChange) % int64(len(s.images)))
	return s.images[nFrame]
}
