package config

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/legacy"
	"log"
	"os"
	"path"
)

// always use 16x9 resolutions

// 1280 x 720 (HD)
//1366 x 768 (HD)
//1600 x 900
//1920 x 1080 (Full HD or 1080p)
//2560 x 1440 (QHD or 1440p)
//3840 x 2160 (4K UHD)

//var WindowWidth = 1280
//var WindowHeight = 720

var WindowWidth = 1920
var WindowHeight = 1080

// var WindowWidth = 3840
// var WindowHeight = 2160
//var WindowWidth = 2560
//var WindowHeight = 1440

type UltimaVConfiguration struct {
	DataFilePath string
	RawDataOvl   []byte
}

func NewUltimaVConfiguration(dataFilePath string) *UltimaVConfiguration {
	uc := UltimaVConfiguration{
		DataFilePath: dataFilePath,
	}

	var err error
	uc.RawDataOvl, err = os.ReadFile(path.Join(dataFilePath, legacy.DATA_OVL))
	if err != nil {
		log.Fatal("Ooof, couldn't read DATA.OVL")
	}

	return &uc
}
