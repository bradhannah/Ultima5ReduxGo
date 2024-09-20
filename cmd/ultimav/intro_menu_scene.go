package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type IntroMenuScene struct {
	introSprites  *sprites.IntroSprites
	borderSprites *sprites.BorderSprites
}

type screenDimensions struct {
	x int
	y int
}

//var ScreenDimension = screenDimensions{x: 800, y: 600}

var ScreenDimension = screenDimensions{x: 1024, y: 768}

// Update method for the IntroMenuScene
func (m *IntroMenuScene) Update(game *Game) error {
	// Switch to the gameplay scene on keypress (e.g., pressing "Enter")
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// Replace this with code to switch to the game scene
		fmt.Println("Switching to Game Scene")
		game.currentScene = &GameScene{}
	}
	return nil
}

// Draw method for the IntroMenuScene
func (m *IntroMenuScene) Draw(screen *ebiten.Image) {
	// Ultima V Logo
	opLogo := sprites.GetXSpriteWithPercents(m.introSprites.Ultima16Logo,
		sprites.PercentXBasedPlacement{
			StartPercentX: .10,
			EndPercentX:   .90,
			StartPercentY: .05,
		})
	screen.DrawImage(m.introSprites.Ultima16Logo, opLogo)

	// Fire animation
	fireSprite := m.introSprites.FlameAnimation.GetCurrentSprite()
	opFire := sprites.GetXSpriteWithPercents(fireSprite,
		sprites.PercentXBasedPlacement{
			StartPercentX: .10,
			EndPercentX:   .90,
			StartPercentY: .275,
		})

	screen.DrawImage(fireSprite, opFire)

	// Redux overlay
	opRedux := sprites.GetXSpriteWithPercents(m.introSprites.Ultima16Logo,
		sprites.PercentXBasedPlacement{
			StartPercentX: .3,
			EndPercentX:   1 - .3,
			StartPercentY: .21,
		})

	screen.DrawImage(m.introSprites.ReduxLogo, opRedux)

	menuBorder, menuBorderOp := m.borderSprites.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(
		305,
		sprites.PercentBasedPlacement{
			StartPercentX: .02,
			EndPercentX:   .98,
			StartPercentY: .50,
			EndPercentY:   .95,
		})
	screen.DrawImage(menuBorder, menuBorderOp)

	// Render the main menu
	ebitenutil.DebugPrint(screen, "Main Menu: Press Enter to Start")
}

func CreateIntroMenuScene() *IntroMenuScene {
	return &IntroMenuScene{
		introSprites:  sprites.NewIntroSprites(),
		borderSprites: sprites.NewBorderSprites(),
	}
}
