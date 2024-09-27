package mainscreen

import (
	"fmt"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultima_v_save/game_state"
	"github.com/hajimehoshi/ebiten/v2"
)

type CharacterSummary struct {
	characterSummaryImage [game_state.MAX_CHARACTERS_IN_PARTY]*ebiten.Image
	FullSummaryImage      *ebiten.Image
	spriteSheet           *sprites.SpriteSheet
	ultimaFont            *text.UltimaFont
	output                *text.Output

	characterSpriteDop *ebiten.DrawImageOptions
}

const lineHeightPercent = .075
const perCharacterSummaryWidth = 16 * 6
const perCharacterSummaryHeight = 9 * 6
const lineSpacing = 35
const fontPoint = 22

func NewCharacterSummary(spriteSheet *sprites.SpriteSheet) *CharacterSummary {
	characterSummary := &CharacterSummary{}

	characterSummary.spriteSheet = spriteSheet
	characterSummary.ultimaFont = text.NewUltimaFont(fontPoint)
	characterSummary.output = text.NewOutput(characterSummary.ultimaFont, lineSpacing)

	characterSummary.FullSummaryImage = ebiten.NewImage(perCharacterSummaryWidth, perCharacterSummaryHeight*game_state.MAX_CHARACTERS_IN_PARTY)

	for i := 0; i < len(characterSummary.characterSummaryImage); i++ {
		characterSummary.characterSummaryImage[i] = ebiten.NewImage(perCharacterSummaryWidth, perCharacterSummaryHeight)
	}

	characterSummary.characterSpriteDop = &ebiten.DrawImageOptions{}
	characterSummary.characterSpriteDop.GeoM.Scale(5, 5)

	return characterSummary
}

func GetTranslateXYByPercent(xPercent float64, yPercent float64) (float64, float64) {
	screenWidth, screenHeight := ebiten.WindowSize()

	// get the x start and end values based on the percent
	var xLeft = float64(screenWidth) * xPercent
	var yTop = float64(screenHeight) * yPercent
	return xLeft, yTop
}

func (c *CharacterSummary) Draw(gameState *game_state.GameState, screen *ebiten.Image) {

	for i := 0; i < len(c.characterSummaryImage); i++ {
		// draw onto single summary
		character := gameState.Characters[i]
		textTopYPercent := (float64(i) * .075) + 0.035

		characterPortrait := c.spriteSheet.GetSprite(character.GetKeySpriteIndex())
		dop := ebiten.DrawImageOptions{}

		c.characterSummaryImage[i].DrawImage(characterPortrait, &dop)

		spriteDop := sprites.GetDrawOptionsFromPercentsForWholeScreen(c.characterSummaryImage[i], sprites.PercentBasedPlacement{
			StartPercentX: .77,
			EndPercentX:   .98,
			StartPercentY: (float64(i))*lineHeightPercent + .03,
			EndPercentY:   (float64(i))*lineHeightPercent + .25,
		})
		screen.DrawImage(c.characterSummaryImage[i], spriteDop)

		leftTextDop := ebiten.DrawImageOptions{}
		leftTextX, leftTextY := GetTranslateXYByPercent(0.815, textTopYPercent)
		leftTextDop.GeoM.Translate(leftTextX, leftTextY)

		leftTextOutput := fmt.Sprintf("%s\n%d/%dHP",
			character.GetNameAsString(),
			character.CurrentHp,
			character.MaxHp)
		c.output.DrawText(screen, leftTextOutput, &leftTextDop)

		rightTextDop := ebiten.DrawImageOptions{}
		rightTextX, rightTextY := GetTranslateXYByPercent(0.98, textTopYPercent)
		rightTextDop.GeoM.Translate(rightTextX, rightTextY)

		rightTextOutput := fmt.Sprintf("%s\n%dMP", game_state.CharacterStatuses.GetById(character.Status).FriendlyName, character.CurrentMp)
		c.output.DrawTextRightToLeft(screen, rightTextOutput, &rightTextDop)
	}

}

//func (c *CharacterSummary) drawSingleSummary(summaryImage *ebiten.Image, gameState *game_state.GameState) {
//	characterPortrait := c.spriteSheet.GetSprite(indexes.Avatar)
//	dop := ebiten.DrawImageOptions{}
//
//	summaryImage.DrawImage(characterPortrait, &dop)
//}
