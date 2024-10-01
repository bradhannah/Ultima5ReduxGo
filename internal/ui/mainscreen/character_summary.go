package mainscreen

import (
	"fmt"
	game_state2 "github.com/bradhannah/Ultima5ReduxGo/pkg/game_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
	"github.com/hajimehoshi/ebiten/v2"
)

type CharacterSummary struct {
	characterSummaryImage [game_state2.MAX_CHARACTERS_IN_PARTY]*ebiten.Image
	spriteSheet           *sprites.SpriteSheet
	ultimaFont            *text.UltimaFont
	output                *text.Output

	characterSpriteDop *ebiten.DrawImageOptions
}

const lineHeightPercent = .075
const perCharacterSummaryWidth = 16 * 6
const perCharacterSummaryHeight = 9 * 6

func NewCharacterSummary(spriteSheet *sprites.SpriteSheet) *CharacterSummary {
	characterSummary := &CharacterSummary{}

	characterSummary.spriteSheet = spriteSheet
	characterSummary.ultimaFont = text.NewUltimaFont(fontPoint)
	characterSummary.output = text.NewOutput(characterSummary.ultimaFont, lineSpacing)

	for i := 0; i < len(characterSummary.characterSummaryImage); i++ {
		characterSummary.characterSummaryImage[i] = ebiten.NewImage(perCharacterSummaryWidth, perCharacterSummaryHeight)
	}

	characterSummary.characterSpriteDop = &ebiten.DrawImageOptions{}
	characterSummary.characterSpriteDop.GeoM.Scale(5, 5)

	return characterSummary
}

func (c *CharacterSummary) Draw(gameState *game_state2.GameState, screen *ebiten.Image) {

	for i := 0; i < len(c.characterSummaryImage); i++ {
		// draw onto single summary
		character := gameState.Characters[i]

		if character.PartyStatus != game_state2.InTheParty {
			continue
		}

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
		leftTextX, leftTextY := sprites.GetTranslateXYByPercent(0.815, textTopYPercent)
		leftTextDop.GeoM.Translate(leftTextX, leftTextY)

		leftTextOutput := fmt.Sprintf("%s\n%d/%dHP",
			character.GetNameAsString(),
			character.CurrentHp,
			character.MaxHp)
		c.output.DrawText(screen, leftTextOutput, &leftTextDop)

		rightTextDop := ebiten.DrawImageOptions{}
		rightTextX, rightTextY := sprites.GetTranslateXYByPercent(0.98, textTopYPercent)
		rightTextDop.GeoM.Translate(rightTextX, rightTextY)

		rightTextOutput := fmt.Sprintf("%s\n%dMP", game_state2.CharacterStatuses.GetById(character.Status).FriendlyName, character.CurrentMp)
		c.output.DrawTextRightToLeft(screen, rightTextOutput, &rightTextDop)
	}

}

//func (c *CharacterSummary) drawSingleSummary(summaryImage *ebiten.Image, gameState *game_state.GameState) {
//	characterPortrait := c.spriteSheet.GetSprite(indexes.Avatar)
//	dop := ebiten.DrawImageOptions{}
//
//	summaryImage.DrawImage(characterPortrait, &dop)
//}
