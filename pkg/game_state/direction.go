package game_state

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/ultimav/references"
	"github.com/hajimehoshi/ebiten/v2"
)

type Direction int

const (
	None Direction = iota
	Up
	Down
	Left
	Right
)

func GetKeyAsDirection(key ebiten.Key) Direction {
	switch key {
	case ebiten.KeyUp:
		return Up
	case ebiten.KeyDown:
		return Down
	case ebiten.KeyLeft:
		return Left
	case ebiten.KeyRight:
		return Right
	}
	return None
}

func (d Direction) GetNewPositionInDirection(position *references.Position) *references.Position {
	var newPosition references.Position
	switch d {
	case Up:
		newPosition = position.GetPositionUp()
	case Down:
		newPosition = position.GetPositionDown()
	case Left:
		newPosition = position.GetPositionToLeft()
	case Right:
		newPosition = position.GetPositionToRight()
	default:
		return nil
	}

	return &newPosition
}

func (d Direction) GetDirectionCompassName() string {
	switch d {
	case Up:
		return "North"
	case Down:
		return "South"
	case Left:
		return "West"
	case Right:
		return "East"
	default:
		panic("unhandled default case")
	}
}
