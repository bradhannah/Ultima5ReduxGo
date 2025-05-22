package references

type Direction int

const (
	NoneDirection Direction = iota
	Up
	Down
	Left
	Right
)

func (d Direction) GetNewPositionInDirection(position *Position) *Position {
	var newPosition *Position
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

	return newPosition
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

func (d Direction) GetOppositeDirection() Direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	case NoneDirection:
		return NoneDirection
	default:
		panic("unhandled default case")
	}
}
