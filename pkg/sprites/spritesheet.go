package sprites

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type SpriteSheet struct {
	SpriteImage *ebiten.Image
}

const (
	origTilesWidth  = 512
	origTilesHeight = 256
	TileSize        = 16
	spritesPerRow   = 32
	spriteRows      = 16
	totalSprites    = spriteRows * spritesPerRow
)

//go:embed assets/orig_u5_tiles.png
var origTiles []byte

func NewSpriteSheet() *SpriteSheet {
	spriteSheet := &SpriteSheet{}

	spriteSheet.SpriteImage = NewPngSprite(origTiles) //ebiten.NewImage(origTilesWidth, origTilesHeight)

	return spriteSheet
}

func (s *SpriteSheet) getSpriteImageRectangle(nSprite int) image.Rectangle {
	topLeft := image.Point{X: (nSprite % spritesPerRow) * TileSize, Y: (nSprite / spritesPerRow) * TileSize}
	bottomRight := image.Point{
		X: topLeft.X + TileSize,
		Y: topLeft.Y + TileSize,
	}
	return image.Rectangle{
		Min: topLeft,
		Max: bottomRight,
	}
}

func (s *SpriteSheet) GetSprite(nSprite int) *ebiten.Image {
	sprite := ebiten.NewImageFromImage(s.SpriteImage.SubImage(s.getSpriteImageRectangle(nSprite)))
	return sprite
}