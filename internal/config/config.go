package config

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/spf13/viper"

	"github.com/bradhannah/Ultima5ReduxGo/internal/files"
	"github.com/bradhannah/Ultima5ReduxGo/pkg/helpers"
)

// always use 16x9 resolutions

// 1280 x 720 (HD)
// 1366 x 768 (HD)
// 1600 x 900
// 1920 x 1080 (Full HD or 1080p)
// 2560 x 1440 (QHD or 1440p)
// 3840 x 2160 (4K UHD)

// var WindowWidth = 1280
// var WindowHeight = 720

// var WindowWidth = 1920
// var WindowHeight = 1080

// var WindowWidth = 3840
// var WindowHeight = 2160
// var WindowWidth = 2560
// var WindowHeight = 1440

type UltimaVConfigurationFlags struct {
	Resolution   int
	FullScreen   bool
	DataFilePath string
}

type UltimaVConfiguration struct {
	// DataFilePath     string
	RawDataOvl       []byte
	allWindowConfigs []ScreenResolution

	SavedConfigData *UltimaVConfigurationFlags
}

func NewUltimaVConfiguration() *UltimaVConfiguration {
	uc := UltimaVConfiguration{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.ultima_v_redux")

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err = viper.SafeWriteConfig()
			if err != nil {
				log.Fatalf("Error writing config file, %s", err)
			}
		} else {
			log.Fatalf("Error reading config file, %s", err)
		}
	}

	uc.SavedConfigData = &UltimaVConfigurationFlags{}
	_ = viper.Unmarshal(&uc.SavedConfigData)
	dfp := viper.GetString("DataFilePath")
	if dfp == "" {
		viper.Set("DataFilePath", "/Users/bradhannah/games/Ultima_5/Gold")
		uc.SavedConfigData.DataFilePath = dfp
	}

	uc.UpdateSaveFile()

	var err error
	uc.RawDataOvl, err = os.ReadFile(path.Join(uc.SavedConfigData.DataFilePath, files.DATA_OVL))
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
	return path.Join(uc.SavedConfigData.DataFilePath, files.LOOK2_DAT)
}

func (uc *UltimaVConfiguration) GetCurrentTrackedWindowResolution() ScreenResolution {
	if ebiten.IsFullscreen() {
		return GetWindowResolutionFromEbiten()
	}
	return uc.allWindowConfigs[uc.SavedConfigData.Resolution]
}

func (uc *UltimaVConfiguration) IncrementHigherResolution() {
	uc.SavedConfigData.Resolution = helpers.Min(len(uc.allWindowConfigs)-1, uc.SavedConfigData.Resolution+1)
	uc.UpdateSaveFile()
}

func (uc *UltimaVConfiguration) DecrementLowerResolution() {
	uc.SavedConfigData.Resolution = helpers.Max(uc.SavedConfigData.Resolution-1, 0)
	uc.UpdateSaveFile()
}

func (uc *UltimaVConfiguration) SetFullScreen(bFullScreen bool) {
	ebiten.SetFullscreen(bFullScreen)
	uc.SavedConfigData.FullScreen = bFullScreen
	uc.UpdateSaveFile()
}

func (uc *UltimaVConfiguration) UpdateSaveFile() {
	viper.Set("Resolution", uc.SavedConfigData.Resolution)
	viper.Set("FullScreen", uc.SavedConfigData.FullScreen)
	_ = viper.WriteConfig()
}

func (uc *UltimaVConfiguration) GetAllNpcFilePaths() []string {
	return []string{
		path.Join(uc.SavedConfigData.DataFilePath, files.TOWNE_NPC),
		path.Join(uc.SavedConfigData.DataFilePath, files.DWELLING_NPC),
		path.Join(uc.SavedConfigData.DataFilePath, files.CASTLE_NPC),
		path.Join(uc.SavedConfigData.DataFilePath, files.KEEP_NPC),
	}
}
