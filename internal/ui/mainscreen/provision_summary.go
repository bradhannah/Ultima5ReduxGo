package mainscreen

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites/indexes"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
)

type ProvisionSummary struct {
	spriteSheet *sprites.SpriteSheet
	ultimaFont  *text.UltimaFont
	output      *text.Output
}

const leftImageStartX = .755

func NewProvisionSummary(spriteSheet *sprites.SpriteSheet) *ProvisionSummary {
	provisionSummary := ProvisionSummary{}
	provisionSummary.spriteSheet = spriteSheet

	provisionSummary.ultimaFont = text.NewUltimaFont(fontPoint)
	provisionSummary.output = text.NewOutput(provisionSummary.ultimaFont, lineSpacing, 1, maxCharsPerLine)

	//provisionSummary.FullSummaryImage = ebiten.NewImage(perCharacterSummaryWidth, perCharacterSummaryHeight*game_state.MAX_CHARACTERS_IN_PARTY)

	return &provisionSummary
}

func (p *ProvisionSummary) Draw(gameState *game_state.GameState, screen *ebiten.Image) {
	//textTopYPercent := .845

	p.drawRow(0.81, screen,
		[3]*ebiten.Image{
			p.spriteSheet.GetSprite(indexes.ItemFood),
			p.spriteSheet.GetSprite(indexes.ItemKey),
			p.spriteSheet.GetSprite(indexes.ItemGem),
		},
		[3]string{
			fmt.Sprintf("%d", gameState.Provisions.QtyFood),
			fmt.Sprintf("%d", gameState.Provisions.QtyKeys),
			fmt.Sprintf("%d", gameState.Provisions.QtyGems),
		},
	)

	p.drawRow(0.875, screen,
		[3]*ebiten.Image{
			p.spriteSheet.GetSprite(indexes.ItemTorch),
			p.spriteSheet.GetSprite(indexes.ItemMoney),
			p.spriteSheet.GetSprite(indexes.HolyFloorSymbol),
		},
		[3]string{
			fmt.Sprintf("%d", gameState.Provisions.QtyTorches),
			fmt.Sprintf("%d", gameState.QtyGold),
			fmt.Sprintf("%d", gameState.Karma),
		},
	)

	p.drawBottomRow(.945, screen, gameState)
}

func (p *ProvisionSummary) drawRow(startY float64, screen *ebiten.Image, rowSprites [3]*ebiten.Image, values [3]string) {
	const percentBetweenImageAndText = 0.035

	percentIncreaseByX := .23 / 4

	// draw top row icons
	const imageOffsetPercent = -0.01
	for i, sprite := range rowSprites {

		dop := sprites.GetDrawOptionsFromPercentsForWholeScreen(sprite,
			sprites.PercentBasedPlacement{
				StartPercentX: leftImageStartX + percentIncreaseByX*float64(i+1) + imageOffsetPercent,
				EndPercentX:   leftImageStartX + percentIncreaseByX*float64(i+1) + 0.02 + imageOffsetPercent,
				StartPercentY: startY,
				EndPercentY:   startY + percentBetweenImageAndText,
			})
		screen.DrawImage(sprite, dop)

		textDop := ebiten.DrawImageOptions{}
		textDop.GeoM.Translate(sprites.GetTranslateXYByPercent(leftImageStartX+percentIncreaseByX*(float64(i)+1), startY+percentBetweenImageAndText+0.005))
		p.output.DrawTextCenter(screen, values[i], &textDop)
	}
}

func (p *ProvisionSummary) drawBottomRow(startY float64, screen *ebiten.Image, state *game_state.GameState) {
	percentIncreaseByX := .23 / 4

	textDop := ebiten.DrawImageOptions{}
	textDop.GeoM.Translate(sprites.GetTranslateXYByPercent(
		leftImageStartX+percentIncreaseByX, startY))
	//p.output.DrawTextCenter(screen, "4-6-139", &textDop)
	p.output.DrawTextCenter(screen, state.DateTime.GetDateAsString(), &textDop)

	textDop.GeoM.Reset()
	textDop.GeoM.Translate(sprites.GetTranslateXYByPercent(
		leftImageStartX+percentIncreaseByX*3, startY))
	p.output.DrawTextCenter(screen, state.DateTime.GetTimeAsString(), &textDop)
}
