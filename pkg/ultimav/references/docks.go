package references

import (
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

type DockReferences []DockReference

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
		*dockRefs = append(*dockRefs, dock)
	}
	return dockRefs
}
