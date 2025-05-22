package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
)

type gameBorders struct {
	outsideBorder                 *ebiten.Image
	outsideBorderOp               *ebiten.DrawImageOptions
	rightSideBorder               *ebiten.Image
	rightSideBorderOp             *ebiten.DrawImageOptions
	topRightCharacterBorder       *ebiten.Image
	topRightCharacterBorderOp     *ebiten.DrawImageOptions
	midRightTextOutputBorder      *ebiten.Image
	midRightTextOutputBorderOp    *ebiten.DrawImageOptions
	bottomRightProvisionsBorder   *ebiten.Image
	bottomRightProvisionsBorderOp *ebiten.DrawImageOptions
}

func (g *GameScene) initializeBorders() {
	mainBorder := sprites.NewBorderSprites()
	g.borders.outsideBorder, g.borders.outsideBorderOp = mainBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, sprites.PercentBasedPlacement{
		StartPercentX: 0,
		EndPercentX:   .75 + .01,
		StartPercentY: 0,
		EndPercentY:   1,
	})

	rightSideBorder := sprites.NewBorderSprites()
	g.borders.rightSideBorder, g.borders.rightSideBorderOp = rightSideBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, sprites.PercentBasedPlacement{
		StartPercentX: .75,
		EndPercentX:   1,
		StartPercentY: 0,
		EndPercentY:   1,
	})

	topRightCharacterBorder := sprites.NewBorderSprites()
	g.borders.topRightCharacterBorder, g.borders.topRightCharacterBorderOp = topRightCharacterBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, sprites.PercentBasedPlacement{
		StartPercentX: .751,
		EndPercentX:   1,
		StartPercentY: 0.00,
		EndPercentY:   0.5,
	})

	midRightTextOutputBorder := sprites.NewBorderSprites()
	g.borders.midRightTextOutputBorder, g.borders.midRightTextOutputBorderOp = midRightTextOutputBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, sprites.PercentBasedPlacement{
		StartPercentX: .751,
		EndPercentX:   1,
		StartPercentY: .485,
		EndPercentY:   0.8,
	})

	bottomRightTextOutputBorder := sprites.NewBorderSprites()
	g.borders.bottomRightProvisionsBorder, g.borders.bottomRightProvisionsBorderOp = bottomRightTextOutputBorder.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(borderWidthScaling, sprites.PercentBasedPlacement{
		StartPercentX: .751,
		EndPercentX:   1,
		StartPercentY: 0.8 - .015,
		EndPercentY:   1,
	})
}
