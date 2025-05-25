package references

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/files"
)

type TalkReferences struct {
	WordDict                *WordDict
	talkDataForSmallMapType map[SmallMapMasterTypes]TalkBytesForSmallMapType
}

func NewTalkReferences(config *config.UltimaVConfiguration, dataOvl *DataOvl) *TalkReferences {
	talkReferences := &TalkReferences{} //nolint:exhaustruct
	talkReferences.talkDataForSmallMapType = make(map[SmallMapMasterTypes]TalkBytesForSmallMapType)

	talkReferences.WordDict = NewWordDict(dataOvl.CompressedWords)

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
		talkReferences.talkDataForSmallMapType[smallMapMasterType], err = LoadFile(talkFile)

		if err != nil {
			log.Fatalf("error loading talk file %s: %v", talkFile, err)
		}
	}

	return talkReferences
}
