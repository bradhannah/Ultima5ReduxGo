package references

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/files"
)

type talkDataForSmallMapType map[SmallMapMasterTypes]TalkBytesForSmallMapType

type TalkReferences struct {
	WordDict                *WordDict
	talkDataForSmallMapType talkDataForSmallMapType
	talkScripts             map[SmallMapMasterTypes][]*TalkScript
}

func NewTalkReferences(config *config.UltimaVConfiguration, dataOvl *DataOvl) *TalkReferences {
	talkReferences := &TalkReferences{} //nolint:exhaustruct
	talkReferences.talkDataForSmallMapType = make(map[SmallMapMasterTypes]TalkBytesForSmallMapType)

	talkReferences.WordDict = NewWordDict(dataOvl.CompressedWords)
	talkReferences.talkDataForSmallMapType = createTalkDataForSmallMapType(config)

	//var oof map[SmallMapMasterTypes][]*TalkScript
	talkScripts := make(map[SmallMapMasterTypes][]*TalkScript)

	for _, smt := range getTalkLocationByFiles() {
		for _, oof := range talkReferences.talkDataForSmallMapType[smt] {
			specificSmallMapTalkData := oof
			script, err := ParseNPCBlob(specificSmallMapTalkData, talkReferences.WordDict)

			if err != nil {
				log.Fatalf("error parsing talk data for %v: %v", smt, err)
			}

			talkScripts[smt] = append(talkScripts[smt], script)
		}
	}

	talkReferences.talkScripts = talkScripts

	return talkReferences
}

func createTalkDataForSmallMapType(config *config.UltimaVConfiguration) talkDataForSmallMapType {
	talkByMap := make(talkDataForSmallMapType)

	for smallMapMasterType := Castle; smallMapMasterType <= Keep; smallMapMasterType++ {
		var talkFile string

		switch smallMapMasterType { //nolint:exhaustive
		case Castle:
			talkFile = files.CASTLE_TLK
		case Keep:
			talkFile = files.KEEP_TLK
		case Towne:
			talkFile = files.TOWNE_TLK
		case Dwelling:
			talkFile = files.DWELLING_TLK
		default:
			log.Fatalf("unhandled default case for small map type %v", smallMapMasterType)
		}

		talkFile = config.GetFileWithFullPath(talkFile)

		var err error
		talkByMap[smallMapMasterType], err = LoadFile(talkFile)

		if err != nil {
			log.Fatalf("error loading talk file %s: %v", talkFile, err)
		}
	}

	return talkByMap
}

func getTalkLocationByFiles() []SmallMapMasterTypes {
	return []SmallMapMasterTypes{Castle, Keep, Towne, Dwelling}
}
