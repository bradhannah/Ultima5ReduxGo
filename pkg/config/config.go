package config

import (
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/legacy"
	"github.com/hajimehoshi/ebiten/v2"
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

//var WindowWidth = 1920
//var WindowHeight = 1080

// var WindowWidth = 3840
// var WindowHeight = 2160
//var WindowWidth = 2560
//var WindowHeight = 1440

type UltimaVConfiguration struct {
	DataFilePath        string
	RawDataOvl          []byte
	allWindowConfigs    []ScreenResolution
	currentWindowConfig int
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

	uc.allWindowConfigs = make([]ScreenResolution, 0)
	uc.allWindowConfigs = append(uc.allWindowConfigs, ScreenResolution{X: 1280, Y: 720})
	uc.allWindowConfigs = append(uc.allWindowConfigs, ScreenResolution{X: 1366, Y: 768})
	uc.allWindowConfigs = append(uc.allWindowConfigs, ScreenResolution{X: 1920, Y: 1080})
	uc.allWindowConfigs = append(uc.allWindowConfigs, ScreenResolution{X: 2560, Y: 1440})
	uc.allWindowConfigs = append(uc.allWindowConfigs, ScreenResolution{X: 3840, Y: 2160})

	return &uc
}

func (uc *UltimaVConfiguration) GetLookDataFilePath() string {
	return path.Join(uc.DataFilePath, legacy.LOOK2_DAT)
}

func (uc *UltimaVConfiguration) GetCurrentTrackedWindowResolution() ScreenResolution {
	if ebiten.IsFullscreen() {
		return GetWindowResolutionFromEbiten()
	}
	return uc.allWindowConfigs[uc.currentWindowConfig]
}

func (uc *UltimaVConfiguration) IncrementHigherResolution() {
	uc.currentWindowConfig = helpers.Min(len(uc.allWindowConfigs)-1, uc.currentWindowConfig+1)
}

func (uc *UltimaVConfiguration) DecrementLowerResolution() {
	uc.currentWindowConfig = helpers.Max(uc.currentWindowConfig-1, 0)
}

func (uc *UltimaVConfiguration) SetFullScreen(bFullScreen bool) {
	ebiten.SetFullscreen(bFullScreen)

}
