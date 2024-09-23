package sprites

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteSheet struct {
	SpriteImage *ebiten.Image
}

const (
	origTilesWidth  = 512
	origTilesHeight = 256
	origTileSize    = 16
)

//go:embed assets/orig_u5_tiles.png
var origTiles []byte

func NewSpriteSheet() *SpriteSheet {
	spriteSheet := &SpriteSheet{}

	spriteSheet.SpriteImage = NewPngSprite(origTiles) //ebiten.NewImage(origTilesWidth, origTilesHeight)

	return spriteSheet
}
