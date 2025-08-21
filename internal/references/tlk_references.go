package references

import (
	"log"

	"github.com/bradhannah/Ultima5ReduxGo/internal/config"
	"github.com/bradhannah/Ultima5ReduxGo/internal/files"
)

type talkDataForSmallMapType map[SmallMapMasterTypes]TalkBytesForSmallMapType

type TalkReferences struct {
	WordDict                *WordDict                             `json:"word_dict" yaml:"word_dict"`
	TalkDataForSmallMapType talkDataForSmallMapType               `json:"talk_data_for_small_map_type" yaml:"talk_data_for_small_map_type"`
	TalkScripts             map[SmallMapMasterTypes][]*TalkScript `json:"talk_scripts" yaml:"talk_scripts"`
}

func NewTalkReferences(config *config.UltimaVConfiguration, dataOvl *DataOvl) *TalkReferences {
	talkReferences := &TalkReferences{} //nolint:exhaustruct
	talkReferences.TalkDataForSmallMapType = make(map[SmallMapMasterTypes]TalkBytesForSmallMapType)

	talkReferences.WordDict = NewWordDict(dataOvl.CompressedWords)
	talkReferences.TalkDataForSmallMapType = createTalkDataForSmallMapType(config)

	talkScripts := make(map[SmallMapMasterTypes][]*TalkScript)

	for _, smt := range getTalkLocationByFiles() {
		// required to make an assumption it starts at 1, and they are all in order - which they are
		for nTalk := 1; nTalk <= len(talkReferences.TalkDataForSmallMapType[smt]); nTalk++ {
			specificSmallMapTalkData := talkReferences.TalkDataForSmallMapType[smt][nTalk]
			script, err := ParseNPCBlob(specificSmallMapTalkData, talkReferences.WordDict)

			if err != nil {
				log.Fatalf("error parsing talk data for %v: %v", smt, err)
			}

			talkScripts[smt] = append(talkScripts[smt], script)
		}
	}

	talkReferences.TalkScripts = talkScripts

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

func (t *TalkReferences) GetTalkScript(smallMapType SmallMapMasterTypes) []*TalkScript {
	return t.TalkScripts[smallMapType]
}

func (t *TalkReferences) GetTalkScriptByNpcIndex(smallMapType SmallMapMasterTypes, npcIndex int) *TalkScript {
	return t.TalkScripts[smallMapType][npcIndex]
}

func (tr *TalkReferences) GetTalkScripts() map[SmallMapMasterTypes][]*TalkScript {
	return tr.TalkScripts
}
