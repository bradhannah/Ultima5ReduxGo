package mainscreen

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/bradhannah/Ultima5ReduxGo/internal/party_state"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/sprites"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/text"
)

type CharacterSummary struct {
	characterSummaryImage [party_state.MAX_CHARACTERS_IN_PARTY]*ebiten.Image
	spriteSheet           *sprites.SpriteSheet
	ultimaFont            *text.UltimaFont
	output                *text.Output

	characterSpriteDop *ebiten.DrawImageOptions
}

const lineHeightPercent = .075
const perCharacterSummaryWidth = 16 * 6
const perCharacterSummaryHeight = 9 * 6
const maxCharsPerLine = 16

func NewCharacterSummary(spriteSheet *sprites.SpriteSheet) *CharacterSummary {
	characterSummary := &CharacterSummary{}

	characterSummary.spriteSheet = spriteSheet
	characterSummary.ultimaFont = text.NewUltimaFont(text.GetScaledNumberToResolution(fontPoint))
	characterSummary.output = text.NewOutput(
		characterSummary.ultimaFont,
		text.GetScaledNumberToResolution(lineSpacing),
		10,
		maxCharsPerLine)

	for i := 0; i < len(characterSummary.characterSummaryImage); i++ {
		characterSummary.characterSummaryImage[i] = ebiten.NewImage(perCharacterSummaryWidth, perCharacterSummaryHeight)
	}

	characterSummary.characterSpriteDop = &ebiten.DrawImageOptions{}
	characterSummary.characterSpriteDop.GeoM.Scale(5, 5)

	return characterSummary
}

func (c *CharacterSummary) Draw(partyState *party_state.PartyState, screen *ebiten.Image) {

	for i := 0; i < len(c.characterSummaryImage); i++ {
		// draw onto single summary
		character := partyState.Characters[i]

		if character.PartyStatus != party_state.InTheParty {
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
		leftTextX, leftTextY := sprites.GetTranslateXYByPercent(sprites.PercentBasedCenterPoint{X: 0.815, Y: textTopYPercent})
		leftTextDop.GeoM.Translate(leftTextX, leftTextY)

		leftTextOutput := fmt.Sprintf("%s\n%d/%dHP",
			character.GetNameAsString(),
			character.CurrentHp,
			character.MaxHp)
		c.output.DrawText(screen, leftTextOutput, &leftTextDop)

		rightTextDop := ebiten.DrawImageOptions{}
		rightTextX, rightTextY := sprites.GetTranslateXYByPercent(sprites.PercentBasedCenterPoint{X: 0.98, Y: textTopYPercent})
		rightTextDop.GeoM.Translate(rightTextX, rightTextY)

		rightTextOutput := fmt.Sprintf("%s\n%dMP", party_state.CharacterStatuses.GetById(character.Status).FriendlyName, character.CurrentMp)
		c.output.DrawTextRightToLeft(screen, rightTextOutput, &rightTextDop)
	}

}

// func (c *CharacterSummary) drawSingleSummary(summaryImage *ebiten.Image, gameState *game_state.GameState) {
//	characterPortrait := c.spriteSheet.GetSprite(indexes.Avatar)
//	dop := ebiten.DrawImageOptions{}
//
//	summaryImage.DrawBorder(characterPortrait, &dop)
// }
