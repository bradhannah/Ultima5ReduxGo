package references

import (
	"strings"

	"github.com/bradhannah/Ultima5ReduxGo/pkg/config"
)

const (
	totalDocks       = 4
	startDockXOffset = 0x4D86
	startDockYOffset = 0x4D8A
)

type DockReference struct {
	Location Location
	Position Position
}

type DockReferences struct {
	docks []DockReference
}

func GetListOfAllLocationsWithDocks() []Location {
	return []Location{Jhelom, Minoc, East_Britanny, Buccaneers_Den}
}

func GetListOfAllLocationsWithDocksAsString() []string {
	return []string{Jhelom.String(), Minoc.String(), East_Britanny.String(), Buccaneers_Den.String()}
}

func NewDocks(gameConfig *config.UltimaVConfiguration) *DockReferences {
	dockRefs := &DockReferences{}
	dockLocations := GetListOfAllLocationsWithDocks()

	for i, loc := range dockLocations {
		dock := DockReference{
			Location: loc,
			Position: Position{
				X: Coordinate(gameConfig.RawDataOvl[startDockXOffset+i]),
				Y: Coordinate(gameConfig.RawDataOvl[startDockYOffset+i]),
			},
		}
		dockRefs.docks = append(dockRefs.docks, dock)
	}
	return dockRefs
}

func (d *DockReferences) GetDockPosition(dockLocation Location) Position {
	for _, dock := range d.docks {
		if dock.Location == dockLocation {
			return dock.Position
		}
	}
	return Position{}
}

func (d *DockReferences) GetDockPositionByString(dockLocation string) Position {
	for _, l := range d.docks {
		if strings.ToLower(l.Location.String()) == strings.ToLower(dockLocation) {
			return l.Position
		}
	}
	return Position{}
}
