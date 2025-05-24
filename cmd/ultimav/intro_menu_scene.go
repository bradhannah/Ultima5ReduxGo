package main

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/internal/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/input"
)

type IntroMenuScene struct {
	introSprites  *sprites.IntroSprites
	borderSprites *sprites.BorderSprites
	ultimaFont    *text.UltimaFont
	keyboard      *input.Keyboard
	config        *config.UltimaVConfiguration

	nCurrentSelection int
}

func (m *IntroMenuScene) InvalidateResolution() {
}

func CreateIntroMenuScene() *IntroMenuScene {
	intro := &IntroMenuScene{
		introSprites:      sprites.NewIntroSprites(),
		borderSprites:     sprites.NewBorderSprites(),
		ultimaFont:        text.NewUltimaFont(24),
		keyboard:          input.NewKeyboard(250),
		nCurrentSelection: 0,
	}
	// todo: get rid of hardcode - obviously
	intro.config = config.NewUltimaVConfiguration()
	if intro.config.SavedConfigData.FullScreen {
		ebiten.SetFullscreen(true)
	}
	return intro
}

func (m *IntroMenuScene) GetUltimaConfiguration() *config.UltimaVConfiguration {
	return m.config
}

var boundKeysIntro = []ebiten.Key{ebiten.KeyDown, ebiten.KeyUp, ebiten.KeyEnter}

// Update method for the IntroMenuScene
func (m *IntroMenuScene) Update(game *Game) error {
	// Switch to the gameplay scene on keypress (e.g., pressing "Enter")

	pressedKey := m.keyboard.GetBoundKeyPressed(&boundKeysIntro)
	if pressedKey == nil {
		return nil
	}

	if !m.keyboard.TryToRegisterKeyPress(*pressedKey) {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// Replace this with code to switch to the game scene
		fmt.Println("Switching to Game Scene")

		game.currentScene = NewGameScene(config.NewUltimaVConfiguration())
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		m.nCurrentSelection = int(math.Max(float64(m.nCurrentSelection)-1, 0))
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		m.nCurrentSelection = int(math.Min(float64(m.nCurrentSelection+1), float64(len(text.IntroChoices)-1)))
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
	const fireStartX = 0.1

	fireSprite := m.introSprites.FlameAnimation.GetCurrentImage()

	//nolint:mnd
	opFire := sprites.GetXSpriteWithPercents(fireSprite.Bounds(),
		sprites.PercentXBasedPlacement{
			StartPercentX: fireStartX,
			EndPercentX:   1 - fireStartX,
			StartPercentY: .35,
		})

	screen.DrawImage(fireSprite, opFire)

	// Redux overlay
	const reduxStartX = .3
	//nolint:mnd
	opRedux := sprites.GetXSpriteWithPercents(m.introSprites.Ultima16Logo.Bounds(),
		sprites.PercentXBasedPlacement{
			StartPercentX: reduxStartX,
			EndPercentX:   1 - reduxStartX,
			StartPercentY: .28,
		})

	screen.DrawImage(m.introSprites.ReduxLogo, opRedux)

	//nolint:mnd
	menuBorder, menuBorderOp := m.borderSprites.VeryPixelatedRoundedBlueWhite.CreateSizedAndScaledBorderSprite(
		400,
		sprites.PercentBasedPlacement{
			StartPercentX: .02,
			EndPercentX:   .98,
			StartPercentY: .61,
			EndPercentY:   .99,
		})
	screen.DrawImage(menuBorder, menuBorderOp)

	m.ultimaFont.DrawIntroChoices(screen, m.nCurrentSelection)
}

// Draw method for the IntroMenuScene
func (m *IntroMenuScene) Draw(screen *ebiten.Image) {
	m.drawStaticGraphics(screen)

	// Render the main menu
	ebitenutil.DebugPrint(screen, "Main Menu: Press Enter to Start")
}
