package references

func (m MissileType) GetStringName() string {
	switch m {
	case MissileNone:
		return "None"
	case MissileArrow:
		return "Arrow"
	case MissileCannonBall:
		return "Cannon Ball"
	case MissileAxe:
		return "Axe"
	case MissileRed:
		return "Red"
	case MissileBlue:
		return "Blue"
	case MissileGreen:
		return "Green"
	case MissileViolet:
		return "Violet"
	case MissileRock:
		return "Rock"
	default:
		return "Unknown"
	}
}

func GetMissileTypeFromString(missileTypeStr string) MissileType {
	switch missileTypeStr {
	case "None":
		return MissileNone
	case "Arrow":
		return MissileArrow
	case "Cannon Ball":
		return MissileCannonBall
	case "Axe":
		return MissileAxe
	case "Red":
		return MissileRed
	case "Blue":
		return MissileBlue
	case "Green":
		return MissileGreen
	case "Violet":
		return MissileViolet
	case "Rock":
		return MissileRock
	default:
		return -1 // Invalid type
	}
}
