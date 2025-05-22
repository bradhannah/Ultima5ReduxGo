package sprites

import (
	_ "embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites/indexes"
)

type SpriteSheet struct {
	SpriteImage *ebiten.Image

	spriteImageCache [1024]*ebiten.Image
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

	spriteSheet.SpriteImage = NewPngSprite(origTiles) // ebiten.NewImage(origTilesWidth, origTilesHeight)

	return spriteSheet
}

func (s *SpriteSheet) getSpriteImageRectangle(spriteIndex indexes.SpriteIndex) image.Rectangle {
	nSprite := int(spriteIndex)
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

func (s *SpriteSheet) GetSprite(nSprite indexes.SpriteIndex) *ebiten.Image {

	if s.spriteImageCache[nSprite] == nil {
		sprite := ebiten.NewImageFromImage(s.SpriteImage.SubImage(s.getSpriteImageRectangle(nSprite)))
		s.spriteImageCache[nSprite] = sprite
		return sprite
	}
	return s.spriteImageCache[nSprite]
}
