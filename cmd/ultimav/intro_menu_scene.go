package main

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type IntroMenuScene struct {
	introSprites  *sprites.IntroSprites
	borderSprites *sprites.BorderSprites
	ultimaFont    *text.UltimaFont
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

func (m *IntroMenuScene) drawStaticGraphics(screen *ebiten.Image) {
	// Ultima V Logo
	const logoStartX = 0.05
	opLogo := sprites.GetXSpriteWithPercents(m.introSprites.Ultima16Logo.Bounds(),
		sprites.PercentXBasedPlacement{
			StartPercentX: logoStartX,
			EndPercentX:   1 - logoStartX,
			StartPercentY: .05,
		})
	screen.DrawImage(m.introSprites.Ultima16Logo, opLogo)

	// Fire animation
	const fireStartX = 0.05
	fireSprite := m.introSprites.FlameAnimation.GetCurrentSprite()
	opFire := sprites.GetXSpriteWithPercents(fireSprite.Bounds(),
		sprites.PercentXBasedPlacement{
			StartPercentX: fireStartX,
			EndPercentX:   1 - fireStartX,
			StartPercentY: .35,
		})

	screen.DrawImage(fireSprite, opFire)

	// Redux overlay
	const reduxStartX = .23
	opRedux := sprites.GetXSpriteWithPercents(m.introSprites.Ultima16Logo.Bounds(),
		sprites.PercentXBasedPlacement{
			StartPercentX: reduxStartX,
			EndPercentX:   1 - reduxStartX,
			StartPercentY: .25,
		})

	screen.DrawImage(m.introSprites.ReduxLogo, opRedux)

	menuBorder, menuBorderOp := m.borderSprites.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(
		305,
		sprites.PercentBasedPlacement{
			StartPercentX: .02,
			EndPercentX:   .98,
			StartPercentY: .58,
			EndPercentY:   .95,
		})
	screen.DrawImage(menuBorder, menuBorderOp)

	m.ultimaFont.DrawIntroChoices(screen, 0)
}

// Draw method for the IntroMenuScene
func (m *IntroMenuScene) Draw(screen *ebiten.Image) {
	m.drawStaticGraphics(screen)

	// Render the main menu
	ebitenutil.DebugPrint(screen, "Main Menu: Press Enter to Start")
}

func CreateIntroMenuScene() *IntroMenuScene {
	return &IntroMenuScene{
		introSprites:  sprites.NewIntroSprites(),
		borderSprites: sprites.NewBorderSprites(),
		ultimaFont:    text.NewUltimaFont(),
	}
}
