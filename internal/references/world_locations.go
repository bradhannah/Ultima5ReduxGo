package references

import "github.com/bradhannah/Ultima5ReduxGo/internal/config"

const (
	xCoordsOffset  = 0x1e9a
	yCoordsOffset  = 0x1ec2
	TotalLocations = 40
)

type WorldLocations struct {
	LargeMapLocationPositions map[Location]WorldLocation `json:"large_map_location_positions" yaml:"large_map_location_positions"`
	// shrinePositions
}

type WorldLocation struct {
	Position Position `json:"position" yaml:"position"`
	Location Location `json:"location" yaml:"location"`
}

func NewWorldLocations(gameConfig *config.UltimaVConfiguration) *WorldLocations {
	wl := &WorldLocations{}
	wl.LargeMapLocationPositions = make(map[Location]WorldLocation)

	for i := 0; i < TotalLocations; i++ {
		x := Coordinate(gameConfig.RawDataOvl[xCoordsOffset+i])
		y := Coordinate(gameConfig.RawDataOvl[yCoordsOffset+i])
		wl.LargeMapLocationPositions[Location(i+1)] = WorldLocation{
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
	for _, worldLocation := range wl.LargeMapLocationPositions {
		if worldLocation.Position.Equals(&position) {
			return worldLocation.Location
		}
	}
	return EmptyLocation
}
