package references

import "github.com/bradhannah/Ultima5ReduxGo/pkg/config"

const xCoordsOffset = 0x1e9a
const yCoordsOffset = 0x1ec2
const TotalLocations = 40

type WorldLocations struct {
	largeMapLocationPositions map[Location]WorldLocation
	//shrinePositions
}

type WorldLocation struct {
	Position Position
	Location Location
}

func NewWorldLocations(gameConfig *config.UltimaVConfiguration) *WorldLocations {

	wl := &WorldLocations{}
	wl.largeMapLocationPositions = make(map[Location]WorldLocation)

	for i := 0; i < TotalLocations; i++ {
		x := Coordinate(gameConfig.RawDataOvl[xCoordsOffset+i])
		y := Coordinate(gameConfig.RawDataOvl[yCoordsOffset+i])
		wl.largeMapLocationPositions[Location(i+1)] = WorldLocation{
			Location: Location(i + 1),
			Position: Position{
				X: x,
				Y: y,
			},
		}
	}
	return wl
}

func (wl *WorldLocations) GetLocationByPosition(position Position) Location {
	for _, worldLocation := range wl.largeMapLocationPositions {
		if worldLocation.Position.Equals(position) {
			return worldLocation.Location
		}
	}
	return EmptyLocation
}
